import "./partials/response_errors.js"
import "./partials/sign-in-up.js"

const checkUsername= () => {
    const username = localStorage.getItem('username')
    if (username) {
        const logout = document.querySelector('.header__footer__logout')
        if (logout) {
            logout.textContent = username.slice(0,1).toUpperCase()
        }
    } else {
        fetch(document.location.href)
        .then(response => {
            if (response.headers.get('HX-Current-Username') === null) {
                console.log('no user')
            } else {
                const username = response.headers.get('HX-Current-Username')
                document.querySelector('.header__footer__logout').textContent = username.slice(0,1).toUpperCase()
                localStorage.setItem('username', username)
            }
        })
        .catch(error => console.error('Error:', error))
    }
}

checkUsername()
