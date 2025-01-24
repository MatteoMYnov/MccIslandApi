const modalData = {
    "2011": {
        "name": "Minecon 2011",
        "date": "7 avril 2011, à Las Vegas au États-Unis",
        "description": "Cette cape a été remise aux participants de la Minecon 2011. Elle était attribuée en fonction de leur pseudonyme, marquant la première et dernière fois que cette méthode a été utilisée."
    },
    "2012": {
        "name": "Minecon 2012",
        "date": "18 juillet 2012, à Disneyland Paris en France",
        "description": "Cette cape a été remise aux participants de la Minecon 2012. Les joueurs recevaient un code par e-mail qu'ils devaient entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2013": {
        "name": "Minecon 2013",
        "date": "1er juillet 2013, à Orlando en Floride",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail qu'ils devaient entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2015": {
        "name": "Minecon 2015",
        "date": "4 et 5 juillet 2015, à Londres au Royaume-Uni",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail uniquement après avoir scanné leur billet sur place, qu'ils devaient ensuite entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2016": {
        "name": "Minecon 2016",
        "date": "24 et 25 septembre 2016, à Anaheim au États-Unis",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail uniquement après avoir scanné leur billet sur place, qu'ils devaient ensuite entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "mojangold": {
        "name": "Cape Mojang (classique)",
        "date": "20 decembre 2010 au 7 octobre 2015",
        "description": "Cette cape a été remise aux employés de Mojang Studios entre ces deux dates."
    },
    "mojang": {
        "name": "Cape Mojang",
        "date": "7 octobre 2015 au 26 juillet 2021",
        "description": "Cette cape a été attribuée aux employés de Mojang Studios entre ces deux dates. Sa couleur a été modifiée pour mieux correspondre à la nouvelle identité visuelle de Mojang Studios."
    },
    "mojangstudios": {
        "name": "Cape Mojang Studios",
        "date": "26 juillet 2021 jusqu'à aujourd'hui",
        "description": "Cette cape a été attribuée aux employés de Mojang Studios entre ces deux dates. Son design a été modifié pour mieux refléter le rebranding du logo réalisé par Johan Aronson."
    },
};

// Fonction pour ouvrir le modal et afficher l'image cliquée
function openModal(event) {
    const modal = document.getElementById("myModal");
    const modalImage = document.querySelector(".modal-image");
    const modalName = document.querySelector(".modal-name");
    const modalDate = document.querySelector(".modal-date");
    const modalDesc = document.querySelector(".modal-desc");
    const clickedImage = event.target;  // L'image sur laquelle on a cliqué
    
    // Récupérer l'ID (title) de l'image cliquée
    const imageTitle = clickedImage.getAttribute("title"); // Prend le titre comme identifiant
    
    // Si les informations sont présentes dans le JSON, les insérer dans le modal
    if (modalData[imageTitle]) {
        modalName.textContent = modalData[imageTitle].name;
        modalDate.textContent = modalData[imageTitle].date;
        modalDesc.textContent = modalData[imageTitle].description;
    } else {
        modalName.textContent = "Nom non trouvé";
        modalDate.textContent = "Date non trouvée";
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
