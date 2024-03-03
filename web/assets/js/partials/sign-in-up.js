document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("sign-in-up-form");
    const submitButton = document.getElementById("form-submit");

    if (!form || !submitButton) {
        return;
    }
    form.addEventListener("input", function () {
        // console.log(form.checkValidity());
        submitButton.disabled = !form.checkValidity();
    });
});

window.togglePassword = (span) => {
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
