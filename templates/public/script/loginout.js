/**
 * Submit the form
 * Check input
 * @param url destination of request
 */
function submitForm(url) {
    
  // get and has password
    let passInput = document.querySelector("#password-input")
    let password = passInput.value;
    passInput.value = hash(password);

    // set token value
    let tokenValue = Math.random() * 10_000;
    document.cookie = "token=" + tokenValue + "; path=/";
    let token = document.querySelector("#token");
    token.value = tokenValue;

    // TODO: check input values - check username in backend, but email and password can both be done here, or maybe just do all in backend (password must be checked here)
    // TODO: in checking inputs, make sure can't be too long

    // Send login form
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url); 
    xhr.onload = function(event){ 

        // if login is good
        if (xhr.status === 200) {

          location.reload();


        } else if (xhr.status === 500) {

          document.getElementById("message").innerText = xhr.responseText;
          passInput.value = "";

          // account created - login
        } else if (xhr.status === 201) {
          passInput.value = password;
          submitForm("/login")
        }
    }; 
    
    var formData = new FormData(document.getElementById("login-form")); 
    xhr.send(formData);
    return false;

}

/**
 * Create a login prompt
 */
function loginPrompt() {
    const body = document.querySelector("body");

    let wholePrompt = document.createElement("div");
    wholePrompt.id = "wholePrompt";

    // fog up the background
    let fog = document.createElement("div");
    fog.id = "fog";
    fog.style.opacity = 0.8;
    fog.setAttribute("onclick", "exitPrompt()")
    wholePrompt.appendChild(fog);

    // create prompt
    let prompt = document.createElement("div");
    prompt.id = "prompt";
    

    // Add header
    let header = document.createElement("h1");
    header.innerText = "OutaTime Events";
    prompt.appendChild(header)

    let header2 = document.createElement("h3");
    header2.innerText = "Sign into your account";
    prompt.appendChild(header2)

    // create form
    let form = document.createElement("form");
    form.id = "login-form";
    form.setAttribute("onsubmit", "return submitForm('/login');")
    
    // create inputs

    // email input
    let email = document.createElement("input");
    email.id = "email-input";
    email.setAttribute("name", "email");
    email.setAttribute("type", "email");
    email.setAttribute("required", "");
    email.setAttribute("placeholder", "Email address");
    form.appendChild(email);

    // password input
    let password = document.createElement("input");
    password.id = "password-input";
    password.setAttribute("name", "password");
    password.setAttribute("type", "password");
    password.setAttribute("required", "");
    password.setAttribute("placeholder", "Password");
    form.appendChild(password);

    // token input
    let token = document.createElement("input");
    token.id = "token";
    token.setAttribute("name", "token");
    token.setAttribute("type", "hidden");
    token.setAttribute("required", "");
    token.setAttribute("value", " ");
    form.appendChild(token);

    let message = document.createElement("span");
    message.id = "message";
    form.appendChild(message);
    
    // login input
    let login = document.createElement("input");
    login.setAttribute("type", "submit");
    login.setAttribute("value", "LOGIN");
    login.id = "submit";
    form.appendChild(login);

    

    // add form to prompt
    prompt.appendChild(form)

    // other links
    let forgot = document.createElement("a");
    forgot.innerText = "Forgot password?";
    forgot.setAttribute("href", ""); // TODO set link
    forgot.style.marginBottom = "20px";
    prompt.appendChild(forgot);

    prompt.appendChild(document.createElement("p"));

    let signup = document.createElement("a");
    signup.innerText = "Don't have an account? Register here";
    signup.setAttribute("onclick", "signupPrompt()");
    prompt.appendChild(signup);

    wholePrompt.appendChild(prompt)
    wholePrompt.style.opacity = 1.0;

    body.appendChild(wholePrompt);

}


/**
 * Get rid of the prompt.
 */
function exitPrompt() {

    let wholePrompt = document.querySelector("#wholePrompt");
    let body = document.querySelector("body");
    body.removeChild(wholePrompt);
}




/**
 * hash
 * @returns 
 */
String.prototype.hashCode = function() {
    var hash = 0,
      i, chr;
    if (this.length === 0) return hash;
    for (i = 0; i < this.length; i++) {
      chr = this.charCodeAt(i);
      hash = ((hash << 5) - hash) + chr;
      hash |= 0; // Convert to 32bit integer
    }
    return hash;
  }

function hash(string) {
    return string.hashCode();
}

/**
 * Log out
 */
function logOut() {
  
  let form = document.createElement("form")
  form.setAttribute("method", "POST");
  form.setAttribute("action", "/logout"); 
  document.querySelector("body").appendChild(form);
  form.submit();

}