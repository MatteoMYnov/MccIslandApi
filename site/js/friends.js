function toggleFriends() {
    const friendsBody = document.getElementById("friendsBody");
    if (friendsBody.style.maxHeight) {
        friendsBody.style.maxHeight = null;
    } else {
        friendsBody.style.maxHeight = friendsBody.scrollHeight + "px";
    }
}