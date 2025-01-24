const modalData = {
    "2011": {
        "name": "Minecon 2011",
        "description": "Cette cape a été remise aux participants de la MINECON 2011, à Las Vegas le 7 avril 2011, La cape leurs étaient rémise par apport à leurs nom, c'était la première fois que cette méthode fut utilisé, et aussi la dernière."
    },
    "2012": {
        "name": "Minecon 2012",
        "description": "Cette cape a été remise aux participants de la MINECON 2012, à Disneyland Paris le 18 juillet 2012, les joueurs ont été remit d'un code par mail qu'ils devaient mettre sur le site minecraft.net pour récupéré leur récompense."
    },
    "2013": {
        "name": "Minecon 2013",
        "description": "Cette cape a été remise aux participants de la MINECON 2013, à Oralando en Floride le 1er juillet 2013, les joueurs ont été remit d'un code par mail qu'ils devaient mettre sur le site minecraft.net pour récupéré leur récompense."
    },
};

// Fonction pour ouvrir le modal et afficher l'image cliquée
function openModal(event) {
    const modal = document.getElementById("myModal");
    const modalImage = document.querySelector(".modal-image");
    const modalName = document.querySelector(".modal-name");
    const modalDesc = document.querySelector(".modal-desc");
    const clickedImage = event.target;  // L'image sur laquelle on a cliqué
    
    // Récupérer l'ID (title) de l'image cliquée
    const imageTitle = clickedImage.getAttribute("title"); // Prend le titre comme identifiant
    
    // Si les informations sont présentes dans le JSON, les insérer dans le modal
    if (modalData[imageTitle]) {
        modalName.textContent = modalData[imageTitle].name;
        modalDesc.textContent = modalData[imageTitle].description;
    } else {
        modalName.textContent = "Nom non trouvé";
        modalDesc.textContent = "Description non disponible";
    }

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
    if (event.target === modal) {
        closeModal();
    }
});
