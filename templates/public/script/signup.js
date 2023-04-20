
function signupPrompt() {
    let prompt = document.getElementById("prompt")

    // empty promt
    prompt.innerHTML = "";

    // build signup form

    // Add header
    let header = document.createElement("h1");
    header.innerText = "OutaTime Events";
    prompt.appendChild(header)

    let header2 = document.createElement("h3");
    header2.innerText = "Create your new account";
    prompt.appendChild(header2)

    // create form
    let form = document.createElement("form");
    form.id = "login-form";
    form.setAttribute("onsubmit", "return submitForm('/signup');")
    
    // create inputs

    // email input
    let email = document.createElement("input");
    email.id = "email-input";
    email.setAttribute("name", "email");
    email.setAttribute("type", "email");
    email.setAttribute("required", "");
    email.setAttribute("placeholder", "Email address");
    form.appendChild(email);

    // username input
    let username = document.createElement("input");
    username.id = "username-input";
    username.setAttribute("name", "username");
    username.setAttribute("type", "text");
    username.setAttribute("required", "");
    username.setAttribute("placeholder", "Username");
    form.appendChild(username);

    // url input
    let url = document.createElement("input");
    url.id = "url-input";
    url.setAttribute("name", "url");
    url.setAttribute("type", "text");
    url.setAttribute("required", "");
    url.setAttribute("placeholder", "Public URL");
    form.appendChild(url);

    // password input
    let password = document.createElement("input");
    password.id = "password-input";
    password.setAttribute("name", "password");
    password.setAttribute("type", "text");
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
    login.setAttribute("value", "CREATE ACCOUNT");
    login.id = "submit";
    form.appendChild(login);

    prompt.appendChild(form)
}