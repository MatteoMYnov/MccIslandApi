document.querySelector("form").addEventListener("submit", function() {
    // Afficher l'overlay de chargement lorsque le formulaire est soumis
    document.getElementById("loadingOverlay").style.display = "flex";
    setTimeout(function() {
        document.getElementById("loadingOverlay").style.display = "none";
    }, 8000);
});

function loading() {
    document.getElementById("loadingOverlay").style.display = "flex";
    setTimeout(function() {
        document.getElementById("loadingOverlay").style.display = "none";
    }, 8000);
}

window.loading = loading;