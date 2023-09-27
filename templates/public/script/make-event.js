
/**
 * Display the HTML content to make the event
 */
function displayMaker() {

    // get body
    body = document.querySelector("#body");

    // store body HTML in variable
    maker = "<h1>Make a new event</h1>";

    // start form
    maker += '<form id="event-form" onsubmit="return makeEvent()">';

    // add inputs
    maker += '<label>Event Name: <input name="name" type="text" required></label><br>';
    
    maker += '<br /><label>Start: <input name="start" type="datetime-local" required></label><br>';
    maker += '<label>End: <input name="end" type="datetime-local" required></label><br>';
    
    // TODO: Add Google Maps API for the location picker.
    maker += '<br /><label>Location: <input name="name" type="text" required></label><br>';
    
    maker += '<br /><label>Image: <input name="name" type="file" accept="image/*"></label><br>';
    
    maker += '<br /><label>Description: <textarea name="descr" maxlength="20000" rows="20" cols="100" required></textarea></label><br>';
    maker += '<label>Summary: <textarea name="tldr" maxlength="180" placeholder="Summarize in 180 characters" rows="2" cols="100" required></textarea></label><br>';
    
    // Radio buttons. Only select one. Public selected by default
    maker += '<br /><label title="Anyone can see this event">Public: <input name="view" type="radio" checked required></label><br>';
    maker += '<label title="Only you can see this event">Private: <input name="view" type="radio" required></label><br>';
    maker += '<label title="Only those with the link can see this event">Hidden: <input name="view" type="radio" required></label><br>';
    

    // add hidden token input
    maker += '<input id="token" name="token" type="hidden" required value=" ">'

    // add message text for 500 response text
    maker += '<p id="message"></p>';

    // submit button
    maker += '<br /><input type="submit" value="Create Event">';

    // end form
    maker += '</form>';

    // display
    body.innerHTML = maker;
}

/**
 * Submit form to /create-event to create a new event
 */
function makeEvent() {

    // create token cookie
    let tokenValue = Math.random() * 10_000;
    document.cookie = "token=" + tokenValue + "; path=/";

    // set token input value
    let token = document.querySelector("#token");
    token.value = tokenValue;

    // TODO: input validation

    // send POST to /create-event
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/create-event"); 
    xhr.onload = function(event){ 

        // if event creation is good
        if (xhr.status === 200) {
            
            // location.reload()

            // for testing. When done, new event will either reload the events page or take to the new event page itself.
            message.style.color = "green";
            message.innerText = xhr.responseText;
          

            // if event creation is bad is bad
        } else if (xhr.status === 500) {

            let message = document.querySelector("#message");

            message.style.color = "red";
            message.innerText = xhr.responseText;
          
        }
    }; 

    var formData = new FormData(document.getElementById("event-form")); 
    xhr.send(formData);

    return false;
}