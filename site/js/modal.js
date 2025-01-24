// Fonction pour ouvrir le modal et afficher l'image cliquée
function openModal(event) {
    const modal = document.getElementById("myModal");
    const modalImage = document.querySelector(".modal-image");
    const clickedImage = event.target;  // L'image sur laquelle on a cliqué
    
    // Mettre l'image cliquée dans le modal
    modalImage.src = clickedImage.src;  // Prendre la source de l'image cliquée
    
    // Afficher le modal
    modal.style.display = "flex";
}

// Fonction pour fermer le modal
function closeModal() {
    const modal = document.getElementById("myModal");
    modal.style.display = "none";  // Masquer le modal
}

// Ouvrir le modal lorsque l'on clique sur une image avec un data-id contenant "using-modal"
const modalTriggers = document.querySelectorAll('[data-id*="using-modal"]');
modalTriggers.forEach(trigger => {
    trigger.addEventListener('click', openModal);
});

// Fermer le modal si l'on clique en dehors de la fenêtre du modal
const modal = document.getElementById("myModal");
modal.addEventListener('click', (event) => {
    // Vérifier si l'utilisateur a cliqué à l'extérieur de la fenêtre du modal (et pas sur son contenu)
    if (event.target === modal) {
        closeModal();
    }
});
