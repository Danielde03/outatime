
/**
 * Display the HTML content to make the event
 */
function displayMaker() {

    // get body
    body = document.querySelector("#body");

    // store body HTML in variable
    maker = "<h1>Make a new event</h1>";

    // start form
    maker += '<form action="/create-event" method="post" onsubmit="return createEvent()">';

    // add inputs
    maker += '<label>Event Name: <input name="name" type="text"></label><br>';
    
    maker += '<br /><label>Start: <input name="start" type="datetime-local"></label><br>';
    maker += '<label>End: <input name="end" type="datetime-local"></label><br>';
    
    maker += '<br /><label>Location: <input name="name" type="text"></label><br>';
    
    maker += '<br /><label>Image: <input name="name" type="file" accept="image/*"></label><br>';
    
    maker += '<br /><label>Description: <textarea name="descr" maxlength="20000" rows="20" cols="100"></textarea></label><br>';
    maker += '<label>Summary: <textarea name="tldr" maxlength="180" placeholder="Summarize in 180 characters" rows="2" cols="100"></textarea></label><br>';
    
    // Radio buttons. Only select one. Public selected by default
    maker += '<br /><label title="Anyone can see this event">Public: <input name="view" type="radio" checked></label><br>';
    maker += '<label title="Only you can see this event">Private: <input name="view" type="radio"></label><br>';
    maker += '<label title="Only those with the link can see this event">Hidden: <input name="view" type="radio"></label><br>';
    

    // add hidden token input
    maker += '<input id="token" name="token" type="hidden" required value=" ">'

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
function createEvent() {

    // create token cookie
    let tokenValue = Math.random() * 10_000;
    document.cookie = "token=" + tokenValue + "; path=/";

    // set token iput value
    let token = document.querySelector("#token");
    token.value = tokenValue;

    // TODO: input validation

    // submit form
    submitForm("/create-event")

    return false;
}