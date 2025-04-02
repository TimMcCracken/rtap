

function main()

    print("*****hello from Lua")
    lbl1 = display:newLabel("body", 50, 50, 200, 0, 0, "Local" )
    lbl2 = display:newLabel("body", 50, 300, 200, 0, 0, "America/New_York" )
    lbl3 = display:newLabel("body", 50, 550, 200, 0, 0, "UTC" )

    dc1 = display:newDigitalClock("body", 100, 50, 200, 0, 1, "Local" )
    dc2 = display:newDigitalClock("body", 100, 300, 200, 0, 1, "America/New_York" ) 
    dc3 = display:newDigitalClock("body", 100, 550, 200, 0, 1, "UTC" )
    
    print("***** Goodbye from Lua")

end