
var taskApp = {
    cor_ans: "",

    getTasksList: async function(params) {
        var res = document.createElement("div");

        await fetch("/api/get/Tasks" + params, {
            method: "GET",
            headers: {
                Token: "kkajka"
            }
        })
        .then(resp => {
            return resp.json();
        })
        .then(resp => {
            if (resp == null) return res;
            
            for (let i = 0; i < resp.length; i++) {
                let nl = document.createElement("a");
                nl.className = "tasklink";
                nl.id = resp[i]["id"];
                nl.href = `/tasks/${resp[i]["id"]}`;
                nl.innerText = `задача #${resp[i]["id"]} (тип: ${resp[i]["type"]})`;
                res.append(nl);
            }

            return res;
        })

        return res;
    },

    getTaskObject: async function(id, GetCookie) {
        var res = document.createElement("div");

        let file = await fetch(`/api/get/Tasks?id=${id}`, {
            method: "GET",
            headers: {
                Token: "kkajka",
            }
        })
        .then(resp => resp.json())
        .then(resp => {
            if (resp == null || resp.length < 1) return "";

            let task = resp[0];
            
            let head = document.createElement("h2");
            head.innerText = `задача #${id} (тип: ${task["type"]})`;
            res.append(head);

            this.cor_ans = task["ans"];

            return task["text"];
        })

        await fetch(`/files/${file}`)
        .then(resp => {
            return resp.text();
        })
        .then(resp => {
            console.log(resp);
            let nl = document.createElement("div");
            nl.innerHTML = resp;
            res.append(nl);
        })

        let ans = document.createElement("form");
        
        let inp = document.createElement("input");
        inp.type = "text";
        inp.name = "ans";

        let lbl = document.createElement("label");
        lbl.innerText = "введите ответ";
        lbl.htmlFor = "ans";

        let btn = document.createElement("button");
        btn.type = "submit";
        btn.innerText = "проверить";

        let vrd = document.createElement("p");

        ans.append(document.createElement("br"), inp, lbl, document.createElement("br"), btn, vrd);

        let userId = await GetCookie("userId");

        ans.addEventListener("submit", e => {
            e.preventDefault();

            let d = new FormData(e.target);

            fetch("/api/submit", {
                method: "POST",
                headers: {
                    Token: "kkajka",
                },
                body: JSON.stringify({
                    TaskId: id,
                    UserId: userId,
                    Answer: d.get("ans"),
                }),
            })
            .then(resp => {
                if (resp.ok) return resp.json();
                else return null;
            })
            .then(resp => {
                if (resp == null) return;
                if (resp["verdict"] == 0) vrd.innerHTML = "неправильный ответ";
                else vrd.innerHTML = "правильный ответ";
            })
        })

        res.append(ans);

        return res;
    },
}

var examApp = {
    getExamsList: async function () {
        let res = document.createElement("div");

        await fetch("/api/get/Exam", {
            method: "GET",
            headers: {
                Token: "kkajka",
            },
        })
        .then(resp => {
            if (resp.ok) return resp.json();
            else return null;
        })
        .then(resp => {
            if (resp == null) return;

            resp.map(
                exam => {
                    let nl = document.createElement("a");

                    nl.href = `/exams/${exam["id"]}`;
                    nl.innerText = `пробник #${exam["id"]}`;
                    nl.id = `${exam["id"]}`;
                    nl.className = "tasklink";

                    res.append(nl);
                }
            )

            return;
        })

        return res;
    }
}