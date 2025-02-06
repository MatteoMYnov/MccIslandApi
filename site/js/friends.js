function toggleRoll(bodyId, arrowId) {
    const RollBody = document.getElementById(bodyId);
    const arrow = document.getElementById(arrowId);

    if (RollBody.style.maxHeight) {
        RollBody.style.maxHeight = null;
        arrow.style.transform = "rotate(0deg)"; // Remet la flèche à sa position initiale
    } else {
        RollBody.style.maxHeight = RollBody.scrollHeight + "px";
        arrow.style.transform = "rotate(-90deg)"; // Fait tourner la flèche
    }
}

function searchPlayer(element) {
    var playerName = element.getAttribute('data-playername'); // Récupère le nom du joueur
    var searchForm = document.querySelector('form[action="/menu"]'); // Trouve le formulaire de recherche
    var inputField = searchForm.querySelector('input[name="q"]'); // Trouve le champ de texte du formulaire
    
    // Remplace la valeur du champ de recherche avec le nom du joueur cliqué
    inputField.value = playerName;
    
    // Soumet le formulaire de recherche pour rediriger vers la page de profil du joueur
    searchForm.submit();
}