import GetCookie from "./general.js";


var UserData = {
    getUserData: async function () {
        let res = document.createElement("div");
        res.innerHTML = `
            not logined
            <a href="/login">войти</a>/<a href="/register">зарегистрироваться</a>
        `;

        let val = GetCookie("user");
        // console.log(val);
        if (val.length != 0) {
            let exit = document.createElement("form");
            exit.method = "POST";

            exit.addEventListener("submit", e => {
                e.preventDefault();

                fetch("/exit", {
                    method: "POST",
                })
                .then(resp => {
                    if (resp.ok) {
                        window.location.replace("/");
                    }
                })
            });
            
            let button = document.createElement("button");
            button.type = "submit";
            button.innerHTML = "выйти";

            let lk = document.createElement("a")
            lk.text = "статистика пользователя";
            lk.href = "/stats";

            exit.appendChild(button);

            res.innerHTML = `имя пользователя: ${val}`;
            
            res.append(document.createElement("br"), lk, exit);
        }

        console.log(res);

        return res;
        
    }
}

export default UserData;