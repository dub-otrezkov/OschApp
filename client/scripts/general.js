function GetCookie(cor) {
    let res = "";

    console.log(document.cookie);

    document.cookie.split("; ").map(
        rawc => {
            console.log(rawc);

            let name = rawc.split('=')[0], val = rawc.split('=')[1];

            if (name == cor && val.length != 0) {
                res = val;
            }
        }
    );

    return res;
}

export default GetCookie