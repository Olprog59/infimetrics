document.body.addEventListener('htmx:responseError', function(evt) {
    evt.preventDefault();
    const errors = document.getElementById('errors');
    // faire un switch qui gère les principales erreurs comme les erreurs 4** et 5**
    switch (evt.detail.xhr.status.toString()[0]) {
        case "4":
            errors.innerHTML = evt.detail.xhr.responseText;
            break;
        case "5":
            errors.innerHTML = evt.detail.xhr.responseText;
            break;
        default:
            errors.innerHTML = "An error occurred. Please try again later.";
    }
});

// Sélectionne l'élément à observer
const targetNode = document.getElementById('errors');

// Options de configuration pour l'observateur (observer les changements de contenu)
const config = { childList: true, subtree: true };

let timeout; // Pour gérer le délai

// Callback à exécuter quand une mutation est observée
const callback = function(mutationsList, observer) {
    for(let mutation of mutationsList) {
        if (mutation.type === 'childList') {
            if (mutation.addedNodes.length) {
                // console.log('Contenu ajouté à #errors');
                targetNode.classList.add('show');

                // Annule le timeout précédent s'il existe
                clearTimeout(timeout);

                // Définit un nouveau timeout pour supprimer le contenu après 10 secondes
                timeout = setTimeout(() => {
                    targetNode.textContent = ''; // Supprime le contenu
                    // console.log('Contenu supprimé après 10 secondes');
                    targetNode.classList.remove('show');
                }, 10000); // 10000 millisecondes = 10 secondes

            }
        }
    }
};

// Crée une instance de l'observateur liée au callback
const observer = new MutationObserver(callback);

// Commence à observer l'élément cible avec la configuration donnée
observer.observe(targetNode, config);
