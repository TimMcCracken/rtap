package rtdsms

import (

	"fmt"
	"math"
)




const NOT_NULL 		= 1
const UNIQUE 		= 2
const INDEXED   	= 4
const NO_PERSIST	= 8






// These constants are used internally for various purposes
// and should never be changed after Release 1
const COL_TYPE_UNDEFINED		= 0

const COL_TYPE_BOOL			= 1

const COL_TYPE_I8			= 2
const COL_TYPE_I16			= 3
const COL_TYPE_I32			= 4
const COL_TYPE_I64			= 5

const COL_TYPE_U8			= 6
const COL_TYPE_U16			= 7
const COL_TYPE_U32			= 8
const COL_TYPE_U64			= 9

const COL_TYPE_F32			= 10
const COL_TYPE_F64			= 11
//const COL_TYPE_COMPLEX64	= 12;
//const COL_TYPE_COMPLEX128	= 13;

const COL_TYPE_TIME			= 14
const COL_TYPE_DURATION		= 15


const COL_TYPE_UTF8			= 16	// UTF-8 variable length, with a fixed (padded) length
const COL_TYPE_BYTE			= 17	// similar to VARCHAR, but bytes rather than UTF-8
const COL_TYPE_ASCII		= 18	// similar to VARBYTE, but only allows 0-127

const COL_TYPE_VAR_UTF8		= 19	// UTF-8 variable length, with a max length
const COL_TYPE_VAR_BYTE		= 20	// similar to VARCHAR, but bytes rather than UTF-8
const COL_TYPE_VAR_ASCII	= 21	// similar to VARBYTE, but only allows 0-127

const COL_TYPE_BLOB_UTF8	= 22	// UTF-8 blob
const COL_TYPE_BLOB_BYTE	= 23	// similar to TEXT, but raw bytes
const COL_TYPE_BLOB_ASCII	= 24	// similar to BYTES, but only allows 0-127

const COL_TYPE_LUA			= 25	// lua script (check the syntax on insert or update)
const COL_TYPE_PYTHON		= 26	// FUTURE ??? 

const COL_TYPE_MONEY		= 27	// 

const COL_TYPE_BITARRAY8	= 28	// boolean arrays have built in processing for individual bits
const COL_TYPE_BITARRAY16	= 29
const COL_TYPE_BITARRAY32	= 30
const COL_TYPE_BITARRAY64	= 31

//pub const NAN: f64 = std::f64::NAN;





/* ----------------------------------------------------------------------------
ColumnDescriptor  stores the definitions of a columns.  
-----------------------------------------------------------------------------*/

type ColumnDescriptor struct {
	name			string
	typeName		string
	length			uint64
	typeConstant	uint16
	notNull			bool
	unique			bool
	indexed 		bool
	noPersist 		bool	// doesn't have to be written to disk (calculated or derived)
	config 			bool 	// only needs to be written to disk if changed
	defaultValue	string
	regex			string  // a regex string to validate against
	check			string	// Lua script
	minLength		uint64
	minInt			int64
	maxInt	 		int64
	minUint			uint64
	maxUint	 		uint64
	minFloat		float64
	maxFloat 		float64
}


// ----------------------------------------------------------------------------
// NewColumnDescriptor is variadic, with the variadic type being 'nay'. This
// allows each column type to have a different set of parameters of different
// types.
// 
// Variadic parameter order for various types:
//
// varUTF8	length, minLength, regexp
// intXX	minValue, maxValue, defaultValue, 
// floatXX	minValue, maxValue, defaultValue, 
//
//
// ----------------------------------------------------------------------------


func NewColumnDescriptor(name string, type_name string, options uint64, params ... any ) (*ColumnDescriptor, error) {

	cd := new(ColumnDescriptor)

	cd.name 		= name
	cd.typeName 	= type_name
	cd.minInt	 	= math.MinInt64
	cd.maxInt		= math.MaxInt64
	cd.minUint	 	= 0
	cd.maxUint		= math.MaxUint64
	cd.minFloat		= math.SmallestNonzeroFloat64
	cd.maxFloat 	= math.MaxFloat64

	// TODOL decode options and set the bits
	if options & NOT_NULL > 0 {
		cd.notNull = true;
	}

	if options & UNIQUE > 0 {
		cd.unique = true;
	}

	if options & INDEXED > 0 {
		cd.notNull = true;
		cd.unique = true;
		cd.indexed = true;
	}

	if options & NO_PERSIST > 0 {
		cd.noPersist = true;
	}

	switch type_name {

	case "varUTF8":
		cd.typeConstant = COL_TYPE_VAR_UTF8
		if len(params) == 0 {
			return nil, fmt.Errorf("Not enough parameters provided. Type varUTF8 must specify length")
		}

		// set the length attribute or return an error
		length, ok := params[0].(uint64)
		if ok == false {
			return nil, fmt.Errorf("parameter [length] must be uint64")
		} 
		cd.length = length

		// Check if we have a minumum length and set it if needed
		if len(params) > 1 {
			min_length, ok := params[1].(uint64)
			if ok == false {
				return nil, fmt.Errorf("parameter [length] must be uint64")
			} 
			cd.minLength = min_length
		} 

		// Check if we have a regex and set it if needed
		if len(params) > 2 {
			regex, ok := params[2].(string)
			if ok == false {
				return nil, fmt.Errorf("parameter [length] must be uint64")
			} 
			cd.regex = regex
		} 

		// return an error because we have too many parameters
		if len(params) > 3 {
				return nil, fmt.Errorf("too many parameters for type [%s].", type_name)
		} 

	case "blobUTF8":
		cd.typeConstant = COL_TYPE_BLOB_UTF8

		// return an error because we have too many parameters
		if len(params) > 3 {
				return nil, fmt.Errorf("too many parameters for type [%s].", type_name)
		} 

	case "int64":
		cd.typeConstant = COL_TYPE_I64
		// TODO: process min, max, default
	case "float64":
		cd.typeConstant = COL_TYPE_F64
		// TODO: process min, max, default

	default:
		return nil, fmt.Errorf("type_name [%s] is unknown", type_name)

	}

	return cd, nil
}


func (cd * ColumnDescriptor) Name() string {
	return cd.name
}

func (cd * ColumnDescriptor) TypeName() string {
	return cd.typeName
}

func (cd * ColumnDescriptor) Length() uint64 {
	return cd.length
}

func (cd * ColumnDescriptor) NotNull() bool {
	return cd.notNull
}

func (cd * ColumnDescriptor) Unique() bool {
	return cd.unique
}

func (cd * ColumnDescriptor) Indexed() bool {
	return cd.indexed
}

func (cd * ColumnDescriptor) NoPersist() bool {
	return cd.noPersist 
}

func (cd * ColumnDescriptor) Config() bool {
	return cd.config
}

func (cd * ColumnDescriptor) DefaultValue() string {
	return cd.defaultValue
}

func (cd * ColumnDescriptor) Regex() string {
	return cd.regex
}

func (cd * ColumnDescriptor) Check() string {
	return cd.check
}

func (cd * ColumnDescriptor) MinLength() uint64 {
	return cd.minLength
}

func (cd * ColumnDescriptor) MinInt() int64 {
	return cd.minInt
}

func (cd * ColumnDescriptor) MaxInt() int64 {
	return cd.maxInt
}

func (cd * ColumnDescriptor) MinUint() uint64 {
	return cd.minUint
}

func (cd * ColumnDescriptor) MaxUint() uint64 {
	return cd.maxUint
}

func (cd * ColumnDescriptor) MinFloat() float64 {
	return cd.minFloat
}

func (cd * ColumnDescriptor) MaxFloat() float64 {
	return cd.maxFloat
}



/*
type ColumnDescriptor struct {
	name			string
	typeName		string
	length			uint64
	typeConstant	uint16
	notNull			bool
	unique			bool
	indexed 		bool
	noPersist 		bool	// doesn't have to be written to disk (calculated or derived)
	config 			bool 	// only needs to be written to disk if changed
	defaultValue	string
	regex			string  // a regex string to validate against
	check			string	// Lua script
	minLength		uint64
	minInt			int64
	maxInt	 		int64
	minUint			uint64
	maxUint	 		uint64
	minFloat		float64
	maxFloat 		float64
}
*/