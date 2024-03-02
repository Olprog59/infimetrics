document.body.addEventListener('htmx:responseError', function(evt) {
    evt.preventDefault();
    const errors = document.getElementById('errors');
    // faire un switch qui gère les principales erreurs comme les erreurs 4** et 5**
    switch (evt.detail.xhr.status.toString()[0]) {
        case "4":
            errors.innerHTML = "This page does not exist or you do not have the rights to access it.";
            break;
        case "5":
            errors.innerHTML = "Internal Server Error. Please try again later.";
            break;
        default:
            errors.innerHTML = "An error occurred. Please try again later.";
    }
});

function togglePassword(span) {
    console.log(span)
    const passwordInput = span.previousElementSibling;
    if (passwordInput.type === "password") {
        passwordInput.type = "text";
        span.children[0].src = "/static/icons/visibility_off.png";
    } else {
        passwordInput.type = "password";
        span.children[0].classList.add("fa-eye-slash");
        span.children[0].src = "/static/icons/visibility.png";
    }
}

