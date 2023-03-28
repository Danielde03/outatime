
/**
 * Update a user's page
 * 
 * Send the request and handle responcces from the backend.
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
    
    var formData = new FormData(document.getElementById("updateForm")); 
    xhr.send(formData);

    return false
}