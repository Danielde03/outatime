
/**
 * Update a user's page
 * 
 * Send the request and handle responces from the backend.
 * 
 * @param url destination of request
 */
function updatePage(url) {

    // set token value
    let tokenValue = Math.random() * 10_000;
    document.cookie = "token=" + tokenValue + "; path=/";
    let token = document.querySelector("#token");
    token.value = tokenValue;

    // set value of isPublic
    let isPublic = document.querySelector("#isPublic");
    isPublic.value = isPublic.checked;

    // send POST to /update-page
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url); 
    xhr.onload = function(event){ 

        // if update is good
        if (xhr.status === 200) {
            
            let message = document.querySelector("#message");

            message.style.color = "green";
            message.innerText = "Page updated.";
          

            // if update is bad
        } else if (xhr.status === 500) {

            let message = document.querySelector("#message");

            message.style.color = "red";
            message.innerText = xhr.responseText;
          
        }
    }; 
    
    var formData = new FormData(document.querySelector(".updateForm")); 
    xhr.send(formData);

    return false
}


/**
 * Update a user's account
 * 
 * Send the request and handle responces from the backend.
 * 
 * @param url destination of request
 * @param formId id of the form to send
 */
function updateAccount(url, formId, tokenId) {

    // set token value
    let tokenValue = Math.random() * 10_000;
    document.cookie = "token=" + tokenValue + "; path=/";
    let token = document.getElementById(tokenId);
    token.value = tokenValue;

    // if updating password
    if (formId === "updatePass") {
        let pass = document.getElementById("password");
        let rep = document.getElementById("rep-password");

        // hash both values
        pass.value = hash(pass.value)
        rep.value = hash(rep.value)

    }

    // send POST to url
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url); 
    xhr.onload = function(event){ 

        // if update is good
        if (xhr.status === 200) {
            
            let message = document.querySelector("#"+formId+"> #message");

            // clear password inputs
            [...document.querySelectorAll("#updatePass > label > input")].forEach(e => {
                e.value = "";
            })

            message.style.color = "green";
            message.innerText = "Page updated. Reload page to see changes";
          

            // if update is bad
        } else if (xhr.status === 500) {

            let message = document.querySelector("#"+formId+"> #message");

            // clear password inputs
            [...document.querySelectorAll("#updatePass > label > input")].forEach(e => {
                e.value = "";
            })

            message.style.color = "red";
            message.innerText = xhr.responseText;
          
        }
    }; 
    
    var formData = new FormData(document.getElementById(formId)); 
    xhr.send(formData);

    return false
}