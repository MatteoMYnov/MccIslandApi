const modalData = {
    "2011": {
        "name": "Minecon 2011",
        "date": "Obtentions: 7 avril 2011, à Las Vegas au États-Unis",
        "description": "Cette cape a été remise aux participants de la Minecon 2011. Elle était attribuée en fonction de leur pseudonyme, marquant la première et dernière fois que cette méthode a été utilisée."
    },
    "2012": {
        "name": "Minecon 2012",
        "date": "Obtentions: 18 juillet 2012, à Disneyland Paris en France",
        "description": "Cette cape a été remise aux participants de la Minecon 2012. Les joueurs recevaient un code par e-mail qu'ils devaient entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2013": {
        "name": "Minecon 2013",
        "date": "Obtentions: 1er juillet 2013, à Orlando en Floride",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail qu'ils devaient entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2015": {
        "name": "Minecon 2015",
        "date": "Obtentions: 4 et 5 juillet 2015, à Londres au Royaume-Uni",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail uniquement après avoir scanné leur billet sur place, qu'ils devaient ensuite entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "2016": {
        "name": "Minecon 2016",
        "date": "Obtentions: 24 et 25 septembre 2016, à Anaheim au États-Unis",
        "description": "Cette cape a été remise aux participants de la Minecon 2013. Les joueurs recevaient un code par e-mail uniquement après avoir scanné leur billet sur place, qu'ils devaient ensuite entrer sur le site minecraft.net pour récupérer leur récompense."
    },
    "mojangold": {
        "name": "Mojang Cape (classique)",
        "date": "Obtentions: 20 decembre 2010 au 7 octobre 2015",
        "description": "Cette cape a été remise aux employés de Mojang Studios entre ces deux dates."
    },
    "mojang": {
        "name": "Mojang Cape",
        "date": "Obtentions: 7 octobre 2015 au 26 juillet 2021",
        "description": "Cette cape a été attribuée aux employés de Mojang Studios entre ces deux dates. Sa couleur a été modifiée pour mieux correspondre à la nouvelle identité visuelle de Mojang Studios."
    },
    "mojangstudios": {
        "name": "Mojang Studios Cape",
        "date": "Obtentions: 26 juillet 2021 jusqu'à aujourd'hui",
        "description": "Cette cape a été attribuée aux employés de Mojang Studios entre ces deux dates. Son design a été modifié pour mieux refléter le rebranding du logo réalisé par Johan Aronson."
    },
    "migrator_cape": {
        "name": "Migrator Cape",
        "date": "Obtentions: 9 juillet 2021",
        "description": "Cette cape a été attribuée aux joueurs possédant un compte avant le 1er décembre 2020 et ayant migré leur compte Mojang vers un compte Microsoft."
    },
    "vanilla_cape": {
        "name": "Vanilla Cape",
        "date": "Obtentions: 31 août 2022",
        "description": "Cette cape a été attribuée aux joueurs ayant acheté Minecraft Java Edition et Minecraft Bedrock Edition avant le 6 juin 2022, en utilisant la même adresse e-mail, avant la fusion des deux éditions."
    },
    "15A": {
        "name": "15th Anniversary Cape",
        "date": "Obtentions: 17 mai 2024 au 21 juin 2024",
        "description": "Cette cape a été offerte aux joueurs qui ont cliqué sur un bouton sur la page de l'anniversaire des 15 ans sur minecraft.net entre le 17 mai 2024 et le 21 juin 2024. Le point de terminaison de l'API permettait aux joueurs de l'échanger après le 21 juin, ce qui a été corrigé quatre jours plus tard. Elle était également appelée « Creepy Cape » sur la page principale de minecraft.net."
    },
    "tiktok": {
        "name": "Follower's Cape",
        "date": "Obtentions: 17 mai 2024 au 18 juin 2024",
        "description": "Cette cape a été offerte aux joueurs qui ont échangé le code lors de l'événement du 15e anniversaire de Minecraft. Les utilisateurs peuvent débloquer un code échangeable pour la cape en regardant un livestream Minecraft sur TikTok. En fonction de la région de l'utilisateur, des conditions supplémentaires étaient imposées, comme regarder pendant un certain nombre de minutes ou taper « minecraft » dans le chat en direct. L'événement a commencé le 17 mai 2024 et s'est terminé le 18 juin 2024. Les joueurs qui ont débloqué la cape l'ont d'abord reçue sur l'édition Bedrock, puis l'ont progressivement étendue à l'édition Java en juillet 2024. Tous les codes pour la cape ont été désactivés le 30 juin 2024. Les joueurs vivant dans des régions où TikTok est interdit ne pouvaient pas déverrouiller la cape, à moins qu'un autre joueur ne leur donne un code. Le 12 juin 2024, un joueur sur Java Edition, sandycheeksvore, a été repéré avec la cape, qui a ensuite été révoquée le 15 juin 2024. Le 1er juillet 2024, la cape a été lentement distribuée à tous les joueurs éligibles",
    },
    "twitch":{
        "name": "Purple Heart Cape",
        "date": "Obtentions: 15 mai 2024 au 31 mai 2024",
        "description": "Cette cape a été offerte aux joueurs qui ont échangé le code lors de l'événement du 15e anniversaire de Minecraft. Les utilisateurs pouvaient débloquer un code Twitch Drop pour la cape en regardant n'importe quel streamer Minecraft avec les Drops activés pendant au moins 15 minutes. Le Drop a été lancé le 15 mai 2024 et s'est terminé le 31 mai 2024. Tous les codes pour la cape ont été désactivés le 30 juin 2024. Les joueurs qui ont débloqué la cape l'ont d'abord reçue sur l'édition Bedrock, avant de l'étendre progressivement à l'édition Java en juillet 2024. Le 20 mai 2024, un joueur sur Java Edition, _Kaptanbey0_, a été repéré avec la cape, qui a ensuite été révoquée le 15 juin 2024. Le 1er juillet 2024, la cape a été lentement distribuée à tous les joueurs éligibles."
    }
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
