function toggleRoll(bodyId, arrowId) {
    const RollBody = document.getElementById(bodyId);
    const arrow = document.getElementById(arrowId);

    if (RollBody.style.maxHeight) {
        RollBody.style.maxHeight = null;
        // RollBody.style.padding = "0 15px"; 
        arrow.style.transform = "rotate(0deg)"; // Remet la flèche à sa position initiale
    } else {
        RollBody.style.maxHeight = RollBody.scrollHeight + 300 + "px";
        // RollBody.style.padding = "5px 15px 5px";
        arrow.style.transform = "rotate(-90deg)"; // Fait tourner la flèche
    }
}