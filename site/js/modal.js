// Fonction pour ouvrir le modal
function openModal() {
    const modal = document.getElementById("myModal");
    modal.style.display = "flex"; // Afficher le modal
}

// Fonction pour fermer le modal
function closeModal() {
    const modal = document.getElementById("myModal");
    modal.style.display = "none"; // Masquer le modal
}

// Ouvrir le modal lorsque l'on clique sur une image avec un data-id contenant "temp"
const modalTriggers = document.querySelectorAll('[data-id*="using-modalBLOCKED"]');
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
