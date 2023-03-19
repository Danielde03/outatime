
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

    // send POST to /update-page
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url); 
    xhr.onload = function(event){ 

        // if update is good
        if (xhr.status === 200) {

          

            // if update is bad
        } else if (xhr.status === 500) {

            // test recieving
            alert(xhr.responseText)
          
        }
    }; 
    
    var formData = new FormData(document.getElementById("updateForm")); 
    xhr.send(formData);

    return false
}