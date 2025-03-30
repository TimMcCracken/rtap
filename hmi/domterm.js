


// ----------------------------------------------------------------------------
// Initialize the function map
// ----------------------------------------------------------------------------
var InitializeFunctionMap = function(functionMap){
    functionMap.set("SetDocumentTitle", setDocumentTitle);
    functionMap.set("SetAttributes", setAttributes);
    functionMap.set("AppendElement", appendElement);
    functionMap.set("SetValue", setValue);
    functionMap.set("SetStyle", setStyle);
}

// ----------------------------------------------------------------------------
// SetDocumentTitle()
// ----------------------------------------------------------------------------
var setDocumentTitle = function(msg){

    // Check that the 'title' field exists in the message data
    title = msg.data["title"];
    if (title === null) {
        console.log("Required field [title] not found in message:", msg);
        return
    }   

    document.title = title
}

// ----------------------------------------------------------------------------
// appendElement()
// ----------------------------------------------------------------------------
var appendElement = function(msg){

    // Check that the 'targetID' field exists in the message
    targetID = msg.targetID
    if (targetID === null) {
        console.log("Required field [targetID] not found in message:", msg);
        return
    }   
    
    // Get the target element
    target = document.getElementById(targetID);
    if (target === null) {
        console.log("Target element not found: ", targetID);
        return
    }

    // Check that the 'tag' field exists in the message data
    tag = msg.data["tag"];
    if (tag === null) {
        console.log("Required field [Tag] not found in message:", msg);
    return
    }   

    // Get the tag
    const newElement = document.createElement(tag);
    if (newElement === null) {
        console.log("Invalid tag:", tag );
        return;
    }

    // Loop through the supplied attributes, but skipping 'tag'
    for (const key in msg.data) {
        if (key == "tag") {
            continue;
        }

        if (key in newElement) {
            const value = msg.data[key];
            newElement.setAttribute(key, value);
        } else {
            console.log(`Element has no attribute: ${key}`);    
        }
    }
    target.appendChild(newElement);
}

// ----------------------------------------------------------------------------
//setValue()
// ----------------------------------------------------------------------------
var setValue = function(msg){

    // Check that the 'targetID' field exists in the message
    targetID = msg.targetID
    if (targetID === null) {
        console.log("Required field [targetID] not found in message:", msg );
        return
    }   

    // get the target
    target = document.getElementById(msg.targetID);
    if (target === null) {
        console.log("Target element not found: ", msg.targetID);
        return
    } 

    // get the new value
    value = msg.data["value"];
    if (tag === null) {
        console.log("Value not found in message:", message);
        return;
    }   
    target.value = value;
}


// ----------------------------------------------------------------------------
//setAttributes()
// ----------------------------------------------------------------------------
var setAttributes = function(msg){

    // Check that the 'targetID' field exists in the message
    targetID = msg.targetID
    if (targetID === null) {
        console.log("Required field [targetID] not found in message:", msg );
        return
    }   

    target = document.getElementById(msg.targetID);
    if (target === null) {
        console.log("Target element not found: ", msg.targetID);
        return
    } 
    for (const key in msg.data) {
        if (key in target) {
            const value = msg.data[key];
            target.setAttribute(key, value);
        } else {
            console.log(`Element has no attribute: ${key}`);    
        }
    }
}


// ----------------------------------------------------------------------------
// setStyle()
// ----------------------------------------------------------------------------
var setStyle = function(msg){

    // Check that the 'targetID' field exists in the message
    targetID = msg.targetID
    if (targetID === null) {
        console.log("Required field [targetID] not found in message:", msg );
        return
    }   

    // get the target element
    target = document.getElementById(msg.targetID);
    if (target === null) {
        console.log("Target element not found: ", msg.targetID);
        return
    } 
    // loop through all the style properties
    for (const key in msg.data) {
        if (key in target.style) {
            const value = msg.data[key];
            target.style.setProperty(key, value);
        } else {
            console.log(`Element style has no attribute: ${key}`);    
        }
    }
}




// ----------------------------------------------------------------------------
// sendEvent()
// ----------------------------------------------------------------------------
var sendMouseEvent = function(event){

    const e = new Object();  

    e.type = event.type
 //   e.view = event.view
    e.id = event.target.id
    e.current_id = event.currentTarget.id    
    e.isTrusted = event.isTrusted
    e.button = event.button //0=left, 1=middle, 2=right
    e.altKey = event.altKey
    e.ctrlKey = event.ctrlKey
    e.shiftKey = event.shiftKey
    e.metaKey = event.metaKey


//    console.log(event);
    const jsonMessage = JSON.stringify(e);
    ws.send(jsonMessage);

    console.log("sent mouse event.")
}