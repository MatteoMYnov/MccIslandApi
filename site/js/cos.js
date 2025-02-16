function changeContent(button, contentId) {
    const allContents = document.querySelectorAll('.content');
    const selectedContent = document.getElementById(contentId);

    // Masquer tout le contenu
    allContents.forEach(content => {
        content.classList.add('hidden'); // Ajoute la classe pour masquer
        content.classList.remove('active'); // Enlève la classe active
    });

    // Afficher le contenu sélectionné
    selectedContent.classList.remove('hidden'); // Affiche l'élément
    selectedContent.classList.add('active'); // Ajoute la classe active

    // Gérer la sélection des boutons
    document.querySelectorAll('.cos-selector button').forEach(btn => btn.classList.remove('selected'));
    button.classList.add('selected');
}