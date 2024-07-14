function Process() {
    var d = document.getElementById("userdata");
    d.innerHTML =  `not authorized <a href="/login">войти</a> / <a href="/register">зарегистрироваться</a>`;


    console.log(document.cookie);
    var t = document.cookie.split("; ")
    console.log(t)
    for (i in t) {
        cont = t[i].split('=');
        console.log(cont);
        if (cont[0] == "username" && cont[1].length > 0) {

            console.log("fofofofo");
            console.log(d);
            d.innerHTML = `logined as: ${cont[1]} <form action="/exit" method="post"><button type="submit">выйти</button></form>`
            console.log(d);
        }
    }
}