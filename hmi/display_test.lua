
function main()


    print("*****hello from Lua")
 
   
    lbl1 = display:newLabel("display", 50, 50, 200, 0, 0, "Local", nil )
    lbl2 = display:newLabel("display", 50, 300, 200, 0, 0, "New_York", nil, nil )
    lbl3 = display:newLabel("display", 50, 550, 200, 0, 0, "UTC", nil, nil)

    lbl4 = display:newLabel("display", 250, 50, 200, 0, 0, "Analog Value", nil, nil)


    dc1 = display:newDigitalClock("display", 100, 50, 200, 0, 1, "Local" )
    dc2 = display:newDigitalClock("display", 100, 300, 200, 0, 1, "America/New_York" ) 
    dc3 = display:newDigitalClock("display", 100, 550, 200, 0, 1, "UTC" )

    dc4 = display:newAnalogValue("display", 300, 50, 200, 0, 1, "XYZ" )

   -- dc4.test()
   
    svg1 = display:newSVG("display", 400, 0, 400, 200, 1, "XYZ" )
    circle1 = svg1:newCircle(50, 50, 20, 0, "red", "", 0)
    rect1   = svg1:newRectangle(100, 30, 100, 40, 0, 0, 0, "Yellow", "green", 5)

    styles = {}
    styles.opacity = 0.25
    rect2   = svg1:newRectangle(250, 30, 100, 40, 0, 0, 0, "Yellow", "", 0, nil, styles)


    print("***** Goodbye from Lua")

end

