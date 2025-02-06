// Sélectionner tous les éléments .player-card
const playerCards = document.querySelectorAll('.player-card');

// Ajouter un événement de clic sur chaque carte
playerCards.forEach(card => {
    card.addEventListener('click', function() {
        // Récupérer le nom du joueur à partir de l'attribut data-playeractualname
        const playerName = card.getAttribute('data-playeractualname');
        
        // Rediriger vers la page /menu?q=playerName
        window.location.href = `/menu?q=${playerName}`;
    });
});

document.addEventListener("DOMContentLoaded", function () {
    function getCurrentPageNumber() {
        const match = window.location.pathname.match(/\/classement\/(\d+)/);
        return match ? parseInt(match[1], 10) : 1; // Par défaut, page 1 si non trouvé
    }

    function navigateToPage(pageNumber) {
        if (pageNumber > 0) { // Empêche d'aller en dessous de 1
            window.location.href = `/classement/${pageNumber}`;
        }
    }

    const goLeftButton = document.querySelector('.go-left');
    const goRightButton = document.querySelector('.go-right');

    if (goLeftButton) {
        goLeftButton.addEventListener('click', function () {
            navigateToPage(getCurrentPageNumber() - 1);
        });
    }

    if (goRightButton) {
        goRightButton.addEventListener('click', function () {
            navigateToPage(getCurrentPageNumber() + 1);
        });
    }
});
