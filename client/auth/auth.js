function Login(form) {

    let data = new FormData(form);

    console.log(data.get("login"), data.get("password"));

    fetch("/login", {
        method: "POST",
        headers: {
            token: "kkajka",
        },
        body: JSON.stringify({
            login: data.get("login"),
            password: data.get("password"),
        }),
        
    })
    .then(resp => {
        if (resp.ok) return null;
        else return resp.json();
    })
    .then(r => {
        if (r == null) window.location.replace("/");
        else document.getElementById("response").innerText = r;
    })
}

function Register(form) {

    let data = new FormData(form);

    console.log(data.get("login"), data.get("password"));

    fetch("/register", {
        method: "POST",
        headers: {
            token: "kkajka",
        },
        body: JSON.stringify({
            login: data.get("login"),
            password: data.get("password"),
        }),
        
    })
    .then(resp => {
        if (resp.ok) return null;
        else return resp.json();
    })
    .then(r => {
        if (r == null) window.location.replace("/");
        else document.getElementById("response").innerText = r;
    })
}