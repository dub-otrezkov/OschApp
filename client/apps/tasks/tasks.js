
var taskApp = {
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
                // console.log(resp[i]);


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

    checkTask: async function (taskId, userId, d) {
        return await fetch("/api/submit", {
            method: "POST",
            headers: {
                Token: "kkajka",
            },
            body: JSON.stringify({
                TaskId: parseInt(taskId),
                Answer: d.get("ans"),
                SessionId: -parseInt(userId),
            }),
        })
        .then(resp => {
            if (resp.ok) return resp.json();
            else return null;
        })
    },
    
    getTask: async function(id, GetCookie, checker) {
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

            checker(id, userId, d)
            .then(resp => {
                if (resp == null) return;
                if (resp["verdict"] == 0) vrd.innerHTML = "неправильный ответ";
                else if (resp["verdict"] == 1) vrd.innerHTML = "правильный ответ";
                else if (resp["verdict"] == 2) vrd.innerHTML = "ответ записан";
                else vrd.innerHTML = "сервис временно не работает. повторите позже";
            })
        })

        res.append(ans);

        return res;
    },
}

var examApp = {
    getExamsList: async function () {
        let res = document.createElement("div");

        await fetch("/api/get/Exams", {
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
    },

    showTask: function (id, len) {
        for (let i = 1; i < len + 1; i++) {
            document.getElementById(`sttr${i}`).className = "";
            document.getElementById(`task${i}`).style.display = "none";
        }
        document.getElementById(`sttr${id}`).className = "chosen_btn";
        document.getElementById(`task${id}`).style.display = "block";
    },

    storeSubmit: async function (taskId, userId, d) {
        return await fetch("/api/esubmit", {
            method: "POST",
            headers: {
                Token: "kkajka",
            },
            body: JSON.stringify({
                TaskId: taskId,
                Answer: d.get("ans"),
                SessionId: 0,
            }),
        })
        .then(resp => {
            if (resp.ok) return resp.json();
            else return null;
        })
    },

    getExam: async function (id, GetCookie) {
        let res = document.createElement("div");

        let exam = (await fetch(`/api/get/Tasklist?examId=${id}`, {
            method: "GET",
            headers: {
                Token: "kkajka",
            },
        })
        .then(resp => {
            if (resp.ok) return resp.json();
            else return [{"tasks": ""}];
        }))

        console.log(exam);

        let btns = document.createElement("div");

        res.append(btns);

        for (let i = 0; i < exam.length; i++) {
            let sttr = document.createElement("button");
            sttr.innerText = i + 1;
            sttr.id = `sttr${i + 1}`;

            // sttr.onclick = ;
            sttr.onclick = function() {examApp.showTask(i + 1, exam.length)};

            btns.append(sttr);

            let task = await taskApp.getTask(exam[i]["taskId"], GetCookie, examApp.storeSubmit);

            task.id = `task${i + 1}`;

            res.append(task);
        }
        return res;
    }
}