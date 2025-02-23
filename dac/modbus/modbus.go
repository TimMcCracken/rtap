package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"go.bug.st/serial"
)

const IP_PROTOCOL_TCP uint8 = 0
const IP_PROTOCOL_UDP uint8 = 1

const PROTO_MODBUS_TCP uint8 = 0
const PROTO_MODBUS_RTU uint8 = 1
const PROTO_MODBUS_ASCII uint8 = 2

type modbus_request struct {
	unit_id       uint8
	func_code     uint8
	sub_func      uint16
	start_addr    uint16
	quantity      uint16
	timeout       time.Duration
	response_chan *chan modbus_response
	transaction   uint16
}

type modbus_response struct {
	err            error // always check this first for communication errors,etc.
	unit_id        uint8
	func_code      uint8
	sub_func       uint16
	address        uint16
	quantity       uint16
	exception_code uint8
	exception_msg  string
	byte_count     uint8
	data           []byte
}

type client_channel struct {
	is_serial bool
	//
	host            string
	port            uint16
	ip_protocol     string
	modbus_protocol string
	channel         chan modbus_request
	transaction     uint16
	connection      net.Conn
	stop_flag       bool
	dial_timeout    time.Duration
	// The serial stuff
	serial_param serial.Mode
	serial_port  serial.Port
}

// Not currently supporting Encapsulated Transport/CAN
// Functions supported for ModbusRTU and ModbusTCP
const FC_READ_COILS uint8 = 1
const FC_READ_DISCRETE_INPUTS uint8 = 2
const FC_READ_HOLDING_REGISTERS uint8 = 3
const FC_READ_INPUT_REGISTERS uint8 = 4
const FC_WRITE_SINGLE_COIL uint8 = 5
const FC_WRITE_SINGLE_REGISTER uint8 = 6
const FC_WRITE_MULTIPLE_COILS uint8 = 15
const FC_WRITE_MULTIPLE_REGISTERS uint8 = 16
const FC_READ_FILE_RECORD uint8 = 20
const FC_WRITE_FILE_RECORD uint8 = 21
const FC_MASK_WRITE_REGISTER uint8 = 22
const FC_READ_WRITE_REGISTERS uint8 = 0x23
const FC_READ_FIFO_QUEUE uint8 = 0x24

// Functions supported for ModbusRTU only.
const FC_READ_EXCEPTION uint8 = 0x07
const FC_DIAGNOSTIC uint8 = 0x08
const FC_GET_COMM_EVENT uint8 = 0x011
const FC_GET_COMM_LOG uint8 = 0x012
const FC_REPORT_SERVER_ID uint8 = 0x017
const FC_READ_DEVICE_IDENTIFICATION uint8 = 0x043

// Exception Codes
const EC_ILLEGAL_FUNCTION uint8 = 0x01
const EC_ILLEGAL_DATA_ADDRESS uint8 = 0x02
const EC_ILLEGAL_DATA_VALUE uint8 = 0x03
const EC_SERVER_DEVICE_FAILURE uint8 = 0x04
const EC_ACKNOWLEDGE uint8 = 0x05
const EC_SERVER_DEVICE_BUSY uint8 = 0x06
const EC_MEMORY_PARITY_ERROR uint8 = 0x08
const EC_GATEWAY_PATH_UNAVAILABLE uint8 = 0x0A
const EC_GATEWAY_TARGET_PATH_FAILED_TO_RESPOND uint8 = 0x0B

const MBAP_SIZE uint16 = 7
const CRC_SIZE uint16 = 2

// Determine the 'endian'
func init() {

	//uint16(crcHi)

	fmt.Printf("ENDIAN: %d\n", uint16(0x01)<<8)

}

// ----------------------------------------------------------------------------
// client_channel.Init(), will initialize a channel for either IP or serial.
// If 'mode' == nil, then it will be TCP. If 'mode' is not nil, then a
// serial port will be used. In that case, the host string should be the
// name of the serial port.
//
// TODO: We need to add timers for the serial port. Maybe we use serial.mode
// internally but have our own serial structure that includes timers and
// options. Or maybe timers needs to be a separate slice.
// ----------------------------------------------------------------------------
func (c *client_channel) Init(host string, port uint16, ip_protocol string,
	modbus_protocol string, mode *serial.Mode) error {

	if mode == nil {

		c.stop_flag = false
		c.dial_timeout = 2 * time.Second
		// Parse the IP protocol
		switch ip_protocol {
		case "tcp":
			c.ip_protocol = ip_protocol
		case "udp":
			c.ip_protocol = ip_protocol
		default:
			return fmt.Errorf("error initing channel: [%s] is not a supported IP protocol", ip_protocol)
		}

		// Parse the Modbus protocol
		switch modbus_protocol {
		case "tcp":
			c.modbus_protocol = modbus_protocol
		case "rtu":
			c.modbus_protocol = modbus_protocol
		default:
			return fmt.Errorf("error initing channel: [%s] is not a supported Modbus protocol", modbus_protocol)
		}
		c.host = host
		c.port = port
		c.is_serial = false
	} else {
		// This is a serial channel
		c.serial_param = *mode
		c.host = host
		c.modbus_protocol = "rtu"
		c.is_serial = true
	}

	// Make the channel that the poll requests will be sent to.
	c.channel = make(chan modbus_request)
	return nil
}

// -----------------------------------------------------------------------------
// client_channel.Open()
// -----------------------------------------------------------------------------
func (c *client_channel) Open() error {

	// TODO: TEST CODE GOING AWAY
	/*
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			log.Fatal("No serial ports found!")
		}
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
		}
		// end test code.
	*/
	var err error

	if !c.is_serial {
		target := fmt.Sprintf("%s:%d", c.host, c.port)
		c.connection, err = net.DialTimeout(c.ip_protocol, target, c.dial_timeout)
		if err != nil {
			return fmt.Errorf("error opening network connection: %s", err)
		}
	} else {
		c.serial_port, err = serial.Open(c.host, &c.serial_param)
		if err != nil {
			return fmt.Errorf("error opening serial port: %s", err)
		}
	}

	go c.requestHandler()

	return nil
}

func (c *client_channel) Close() {
	c.stop_flag = true
	if !c.is_serial {
		c.connection.Close()
	} else {
		c.serial_port.Close()
	}
}

// ----------------------------------------------------------------------------
// ReadCoils
// ----------------------------------------------------------------------------
func (c *client_channel) ReadCoils(unit_id uint8, start uint16,
	quantity uint16, timeout time.Duration) ([]byte, error) {

	// parse the variables.
	if quantity > 2000 {
		return nil, fmt.Errorf("specified quantity [%d] is too large", quantity)
	}

	resp_chan := make(chan modbus_response)
	// send the request
	c.poll(unit_id, FC_READ_COILS, 0, start, quantity, timeout, &resp_chan)
	// wait for the response
	response := <-resp_chan

	if response.err != nil {
		return nil, fmt.Errorf("%s", response.err)
	}

	if response.exception_code == 0 {
		return response.data, nil
	} else {
		return nil, fmt.Errorf(response.exception_msg)
	}
}

// ----------------------------------------------------------------------------
// ReadDiscreteInputs
// ----------------------------------------------------------------------------
func (c *client_channel) ReadDiscreteInputs(unit_id uint8, start uint16,
	quantity uint16, timeout time.Duration) ([]byte, error) {

	// parse the variables.
	if quantity > 2000 {
		return nil, fmt.Errorf("specified quantity [%d] is too large", quantity)
	}

	resp_chan := make(chan modbus_response)
	// send the request
	c.poll(unit_id, FC_READ_DISCRETE_INPUTS, 0, start, quantity, timeout, &resp_chan)
	// wait for the response
	response := <-resp_chan

	if response.err != nil {
		return nil, fmt.Errorf("%s", response.err)
	}

	if response.exception_code == 0 {
		return response.data, nil
	} else {
		return nil, fmt.Errorf(response.exception_msg)
	}
}

// ----------------------------------------------------------------------------
// ReadInputRegisters
// ----------------------------------------------------------------------------
func (c *client_channel) ReadInputRegisters(unit_id uint8, start uint16,
	quantity uint16, timeout time.Duration) ([]byte, error) {

	// parse the variables.
	if quantity > 123 {
		return nil, fmt.Errorf("specified quantity [%d] is too large", quantity)
	}

	resp_chan := make(chan modbus_response)
	// send the request
	c.poll(unit_id, FC_READ_INPUT_REGISTERS, 0, start, quantity, timeout, &resp_chan)
	// wait for the response
	response := <-resp_chan

	if response.err != nil {
		return nil, fmt.Errorf("%s", response.err)
	}

	if response.exception_code == 0 {
		return response.data, nil
	} else {
		return nil, fmt.Errorf(response.exception_msg)
	}
}

// ----------------------------------------------------------------------------
// ReadHoldingRegisters
// ----------------------------------------------------------------------------
func (c *client_channel) ReadHoldingRegisters(unit_id uint8, start uint16,
	quantity uint16, timeout time.Duration) ([]byte, error) {

	// parse the variables.
	if quantity > 123 {
		return nil, fmt.Errorf("specified quantity [%d] is too large", quantity)
	}

	resp_chan := make(chan modbus_response)
	// send the request
	c.poll(unit_id, FC_READ_HOLDING_REGISTERS, 0, start, quantity, timeout, &resp_chan)
	// wait for the response
	response := <-resp_chan

	if response.exception_code == 0 {
		return response.data, nil
	} else {
		return nil, fmt.Errorf(response.exception_msg)
	}
}

// ----------------------------------------------------------------------------
// WriteSingleRegister
// ----------------------------------------------------------------------------
func (c *client_channel) WriteSingleRegister(unit_id uint8, address uint16,
	value uint16, timeout time.Duration) ([]byte, error) {

	resp_chan := make(chan modbus_response)
	// send the request
	c.poll(unit_id, FC_WRITE_SINGLE_REGISTER, 0, address, value, timeout, &resp_chan)
	// wait for the response
	response := <-resp_chan

	if response.exception_code == 0 {
		return response.data, nil
	} else {
		return nil, fmt.Errorf(response.exception_msg)
	}
}

// ----------------------------------------------------------------------------
// client_channel.poll() builds a request block, and sends it through a
// GO channel to the requestHandler
// ----------------------------------------------------------------------------
func (c *client_channel) poll(unit_id uint8, func_code uint8, sub_func uint16,
	start uint16, quantity uint16, timeout time.Duration,
	resp_chan *chan modbus_response) {

	req := modbus_request{unit_id, func_code, sub_func, start, quantity, timeout, resp_chan, 0}
	c.channel <- req

}

// ----------------------------------------------------------------------------
// request_handler() is the thread that listens for requests on a GO channel.
// It then generates/sends a request and extracts the data or handles an error.
// ----------------------------------------------------------------------------
func (c *client_channel) requestHandler() {

	// This loop runs forever, stopping only when the shutdown flag is set.
	for {
		/*	if c.stop_flag {
			return
		}*/

		time.Sleep(1 * time.Second) // testing, remove later

		req := <-c.channel
		resp := new(modbus_response) // create a response structure to return

		// the variables we need for processing the data.
		var tx_msg_size uint16
		var rx_msg_size uint16
		//var tx_index int
		var rx_index int

		// Calculate the transmit and receive buffer sizes.
		// First add the overhead bytes which depend on the Modbus protocol used
		if c.modbus_protocol == "tcp" {
			tx_msg_size = MBAP_SIZE
			rx_msg_size = MBAP_SIZE
			// always increment the transaction counter so we can sort out the received data.
			c.transaction++
			req.transaction = c.transaction
		} else {
			tx_msg_size = CRC_SIZE + 1 //add 1 for the unit_id
			rx_msg_size = CRC_SIZE + 1 //add 1 for the unit_id
		}

		// Now add the application protocol sizes
		switch req.func_code {
		case FC_READ_COILS, FC_READ_DISCRETE_INPUTS:
			tx_msg_size += 5 // function code, start address, quantity
			rx_msg_size += 2 // function code, byte count
			if req.quantity%8 > 0 {
				rx_msg_size += (req.quantity / 8) + 1
			} else {
				rx_msg_size += (req.quantity / 8)
			}
		case FC_READ_INPUT_REGISTERS, FC_READ_HOLDING_REGISTERS:
			tx_msg_size += 5                      // function code, start address, quantity
			rx_msg_size += (req.quantity * 2) + 2 // function code, byte count, values
		case FC_WRITE_SINGLE_COIL, FC_WRITE_SINGLE_REGISTER:
			tx_msg_size += 5
			rx_msg_size += 5
		case FC_WRITE_MULTIPLE_COILS, FC_WRITE_MULTIPLE_REGISTERS:
			tx_msg_size += 6
			rx_msg_size += 5
			if req.quantity%8 > 0 {
				tx_msg_size += (req.quantity / 8) + 1
			} else {
				tx_msg_size += (req.quantity / 8)
			}
		default:
			// Invalid function code, log error
			//TODO: LOG AN ERROR
			fmt.Printf("Error in Modbus requestHandler() - unsupported function code: %d.", req.func_code)
			continue // And the next request.
		}
		txBuf := make([]byte, 0, tx_msg_size)
		rxBuf := make([]byte, rx_msg_size)

		if c.modbus_protocol == "tcp" {
			txBuf = append(txBuf, byte((c.transaction&0xFF00)>>8)) //transaction HI
			txBuf = append(txBuf, byte(c.transaction&0xFF))        // transaction LO
			txBuf = append(txBuf, 0)                               // protocol HI
			txBuf = append(txBuf, 0)                               // protocol LO
			txBuf = append(txBuf, 0)                               // Length HI
			txBuf = append(txBuf, byte(tx_msg_size-6))             // Length LO
		}
		txBuf = append(txBuf, req.unit_id)
		txBuf = append(txBuf, req.func_code)
		txBuf = append(txBuf, byte((req.start_addr&0xFF00)>>8))
		txBuf = append(txBuf, byte(req.start_addr&0xFF))
		txBuf = append(txBuf, byte((req.quantity&0xFF00)>>8))
		txBuf = append(txBuf, byte(req.quantity&0xFF))

		// Add the CRC
		if c.modbus_protocol == "rtu" {
			msg := txBuf[0 : tx_msg_size-2]
			crc := crc16(msg)

			// Modbus CRC is sent low byte first
			txBuf = append(txBuf, byte(crc&0x00FF))
			txBuf = append(txBuf, byte((crc&0xFF00)>>8))
		}
		fmt.Printf("TX Dump: %s\n", hex.Dump(txBuf))

		time.Sleep(1 * time.Second)

		if !c.is_serial { // using a network connection
			//fmt.Printf("Timeout: %d\n", 1*time.Second)
			n, err := c.connection.Write(txBuf)
			if err != nil {
				resp.err = fmt.Errorf("error writing network: %s", err)
				*req.response_chan <- *resp
				continue
			}
			if n != int(tx_msg_size) {
				resp.err = fmt.Errorf("error writing network could not write all bytes: %d of %d", n, tx_msg_size)
				*req.response_chan <- *resp
				continue
			}

			fmt.Printf("Tx'd %d bytes\n", n)

			// -----------------------------------------------------------------
			// Read the response
			// -----------------------------------------------------------------
			err = c.connection.SetReadDeadline(time.Now().Add(req.timeout))
			if err != nil {
				resp.err = fmt.Errorf("error setting timeout: %s", err)
				*req.response_chan <- *resp
				continue
			}
			n, err = c.connection.Read(rxBuf)

			fmt.Printf("RX Dump: %d %s", n, hex.Dump(rxBuf))

			if err != nil {
				resp.err = fmt.Errorf("error reading enet: %s", err)
				*req.response_chan <- *resp
				continue
			}
			fmt.Printf("Rec'd %d of %d bytes\n", n, rx_msg_size)
			if n != int(rx_msg_size) {
				resp.err = fmt.Errorf("enet wrong RX message size: got %d expected %d", n, rx_msg_size)
				*req.response_chan <- *resp
				continue
			}

		} else { // Using a serial connection. // TODO: Add modem control
			_ = c.serial_port.ResetInputBuffer()
			_ = c.serial_port.ResetOutputBuffer()
			n, err := c.serial_port.Write(txBuf)
			if err != nil {
				fmt.Printf("Error writing serial: %s\n", err)
				resp.err = err
				*req.response_chan <- *resp
				continue
			}
			if n != int(rx_msg_size) {
				// TODO LOG A MESSAGE
				fmt.Printf("Error writing serial: %s\n", err)
			}

			c.serial_port.SetReadTimeout(2 * time.Second)

			var rx_count uint16
			for rx_count = 0; rx_count < rx_msg_size; {
				fmt.Printf("%d ", rx_count)
				// TO CHECK IS THIS OVERWRITING RATHER THAN APPENDING?
				n, err = c.serial_port.Read(rxBuf)
				rx_count += uint16(n)
			}

			if rx_count != rx_msg_size {
				// TODO LOG A MESSAGE
				fmt.Printf("Error reading serial: %s\n", err)
			}
		}

		// ---------------------------------------------------------------------
		// process the response
		// ---------------------------------------------------------------------
		if c.modbus_protocol == "tcp" {
			resp.data = rxBuf[9 : rx_msg_size-1]
			rx_index = 6
		} else { //modbus_protocol == rtu
			// Check the CRC.
			msg := rxBuf[0 : rx_msg_size-2]

			//	fmt.Printf("RX Msg: %s\n", hex.Dump(msg))

			crc := rxBuf[rx_msg_size-2 : rx_msg_size]
			if binary.LittleEndian.Uint16(crc) != crc16(msg) {
				// TODO LOG A MESSAGE
				fmt.Printf("rx crc error")
			}
			// Check the unit_id
			if msg[0] != req.unit_id {
				// TODO LOG A MESSAGE
				fmt.Printf("rx wrong unit_id")
			}

			resp.data = rxBuf[2 : rx_msg_size-4] // subtract unit_id, func_code and CRC
			rx_index = 0
		}
		resp.unit_id = rxBuf[rx_index]
		rx_index++
		resp.func_code = rxBuf[rx_index]
		rx_index++

		//fmt.Printf("Unit ID: %d  FuncCode: %d\n", resp.unit_id, resp.func_code)
		//	fmt.Printf("[[%x]]", rxBuf)

		// --------------------------------------------------------------------
		// process the rec'd data. First we check for an exception response. If
		// we got an exception response, we log the exception, build the
		// response, and return it to the calling goroutine.
		// TODO: Add logging to each of the exceptions.
		// --------------------------------------------------------------------
		if resp.func_code > 0x80 { // Exception codes = func_code + 0x80.
			resp.exception_code = rxBuf[rx_index]
			switch resp.exception_code {
			case EC_ILLEGAL_FUNCTION:
				resp.exception_msg = fmt.Sprintf("Illegal function: %d", req.func_code)
			case EC_ILLEGAL_DATA_ADDRESS:
				resp.exception_msg = fmt.Sprintf("Illegal data address: %d", req.start_addr)
			case EC_ILLEGAL_DATA_VALUE:
				resp.exception_msg = "Illegal data value"
			case EC_SERVER_DEVICE_FAILURE:
				resp.exception_msg = "Server device failure"
			case EC_ACKNOWLEDGE:
				resp.exception_msg = "Acknowledge"
			case EC_SERVER_DEVICE_BUSY:
				resp.exception_msg = "Server device busy"
			case EC_MEMORY_PARITY_ERROR:
				resp.exception_msg = "Memory parity error"
			case EC_GATEWAY_PATH_UNAVAILABLE:
				resp.exception_msg = "Gateway path unavailable"
			case EC_GATEWAY_TARGET_PATH_FAILED_TO_RESPOND:
				resp.exception_msg = "Gateway target path failed to respond"
			}

			*req.response_chan <- *resp
			continue
		}
		resp.exception_code = 0
		resp.exception_msg = "OK"
		switch resp.func_code {
		case FC_READ_COILS, FC_READ_DISCRETE_INPUTS, FC_READ_HOLDING_REGISTERS, FC_READ_INPUT_REGISTERS:
			resp.byte_count = rxBuf[rx_index]
			rx_index++
		case FC_WRITE_SINGLE_COIL, FC_WRITE_SINGLE_REGISTER:
			resp.address = (uint16(rxBuf[rx_index]) * 2) + uint16(rxBuf[rx_index+1])
			rx_index += 2
		case FC_WRITE_MULTIPLE_COILS, FC_WRITE_MULTIPLE_REGISTERS:
			resp.address = (uint16(rxBuf[rx_index]) * 2) + uint16(rxBuf[rx_index+1])
			rx_index += 2
			resp.quantity = (uint16(rxBuf[rx_index]) * 2) + uint16(rxBuf[rx_index+1])
			rx_index += 2
		}

		// Copy the data to the response.
		resp.data = make([]byte, resp.byte_count)
		n := copy(resp.data, rxBuf[rx_index:rx_index+int(resp.byte_count)])

		// TODO: Do something with n?
		n = n

		//fmt.Printf("%d %d, %d, [%x]\n", n, rx_index, resp.byte_count, rxBuf[rx_index:rx_index+int(resp.byte_count)])
		*req.response_chan <- *resp
	}
}

var channel client_channel

// ----------------------------------------------------------------------------
// Start()
// ----------------------------------------------------------------------------
func Start() {

	fmt.Println("Starting Modbus subsystem.")

	err := channel.Init("192.168.0.253", 502, "tcp", "rtu", nil)
	if err != nil {
		fmt.Printf("Error initing client: %s\n", err)
	} else {
		fmt.Printf("Modbus client inited OK.\n")
	}

	err = channel.Open()
	if err != nil {
		fmt.Printf("Error opening client: %s\n", err)
	} else {
		fmt.Printf("Modbus client opened OK.\n")
	}

	for x := 0; x < 3; x++ {

		_, err := channel.ReadCoils(1, 0, 10, 1*time.Second)

		if err != nil {
			fmt.Printf("ReadDiscreteInputs() error: %s\n", err)
		} else {
			//fmt.Printf("Got back: [%x]\n", data)
		}
	}

	/*
		for x := 0; x < 3; x++ {
			_, err := channel.WriteSingleRegister(1, 0, 3, 1*time.Second)

			if err != nil {
				fmt.Printf("WriteSingle() error: %s\n", err)
			} else {
				//fmt.Printf("Got back: [%x]\n", data)
			}

		}
	*/
}
