// Fonction pour ouvrir le modal et afficher l'image cliquée
function openModal(event) {
    const modal = document.getElementById("myModal");
    const modalImage = document.querySelector(".modal-image");
    const modalName = document.querySelector(".modal-name");
    const modalDate = document.querySelector(".modal-date");
    const modalDesc = document.querySelector(".modal-desc");
    const modalImagesContainer = document.querySelector(".modal-images");
    const clickedImage = event.target;

    // Récupérer le titre (ID) de l'image cliquée
    const utilName  = clickedImage.getAttribute("data-name");
    const modalDataEntry = modalData[utilName ] || {};
    const images = modalDataEntry.images || [];

    // Mettre à jour les informations principales du modal
    modalName.textContent = modalDataEntry.name || "Nom non trouvé";
    modalDate.textContent = modalDataEntry.date || "Date non trouvée";
    modalDesc.textContent = modalDataEntry.description || "Description non disponible";

    // Vider le conteneur des images pour éviter des doublons
    modalImagesContainer.innerHTML = "";

    // Ajouter dynamiquement les images s'il y en a
    if (images.length > 0) {
        images.forEach((image) => {
            const imgElement = document.createElement("img");
            imgElement.src = image.URL || "../img/capes/default-placeholder.png";
            imgElement.alt = image.Alt || "Image";
            imgElement.title = image.Alt || "Image";
            imgElement.className = "modal-image-item";
            imgElement.draggable = "false"
            modalImagesContainer.appendChild(imgElement);
        });
    } else {
        // Si aucune image, cacher le conteneur
        modalImagesContainer.style.display = "none";
    }

    // Mettre l'image cliquée dans le modal (visuel principal)
    modalImage.src = clickedImage.src;

    // Afficher le modal
    modal.style.display = "flex";

    // Désactiver le scroll de l'arrière-plan
    document.body.classList.add("modal-open");
}

// Fonction pour fermer le modal
function closeModal() {
    const modal = document.getElementById("myModal");
    modal.style.display = "none"; // Masquer le modal

    // Réactiver le scroll de l'arrière-plan
    document.body.classList.remove("modal-open");

    // Réinitialiser le conteneur d'images
    const modalImagesContainer = document.querySelector(".modal-images");
    modalImagesContainer.style.display = "block"; // Réactiver l'affichage pour les prochains modals
}

// Ouvrir le modal lorsque l'on clique sur une image avec un data-id contenant "using-modal"
const modalTriggers = document.querySelectorAll('[data-id*="using-modal"]');
modalTriggers.forEach((trigger) => {
    trigger.addEventListener("click", openModal);
});

// Fermer le modal si l'on clique en dehors de la fenêtre du modal
const modal = document.getElementById("myModal");
modal.addEventListener("click", (event) => {
    if (event.target === modal) {
        closeModal();
    }
});
