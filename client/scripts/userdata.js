function getUserData() {
    let res = "";

    document.cookie.split(";").map(
        rawc => {
            console.log(rawc);
        }
    );

}