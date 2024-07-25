
var UserData = {
    getCookie: async function (cor) {
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
    },

    getUserData: async function () {
        let res = document.createElement("div");
        res.innerHTML = `
            not logined
            <a href="/login">войти</a>/<a href="/register">зарегистрироваться</a>
        `;

        let val = await this.getCookie("user");
        // console.log(val);
        if (val.length != 0) {
            let exit = document.createElement("form");
            exit.method = "POST";

            exit.addEventListener("submit", e => {
                e.preventDefault();

                fetch("/exit", {
                    method: "POST",
                    headers: {
                        token: "kkajka",
                    },
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

            exit.appendChild(button);

            res.innerHTML = `имя пользователя: ${val}`;
            res.appendChild(exit);
        }

        console.log(res);

        return res;
        
    }
}