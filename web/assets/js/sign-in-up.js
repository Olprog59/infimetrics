document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("sign-in-up-form");
    const submitButton = document.getElementById("form-submit");

    form.addEventListener("input", function () {
        // console.log(form.checkValidity());
        submitButton.disabled = !form.checkValidity();
    });
});
