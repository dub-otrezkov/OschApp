
var taskApp = {
    getTasksList: async function(params, elemid) {
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
            console.log(resp);
            
            for (let i = 0; i < resp.length; i++) {
                let nl = document.createElement("a");
                nl.className = "tasklink";
                nl.id = resp[i]["id"];
                nl.href = `/tasks/${resp[i]["id"]}`;
                nl.innerText = `задача #${resp[i]["id"]} (тип: ${resp[i]["type"]})`;
                res.append(nl);
            }

            console.log(res);

            return res;
        })

        return res;
    },

    getTaskObject: async function(id) {
        var res = document.createElement("div");

        let file = await fetch(`/api/get/Tasks?id=${id}`, {
            method: "GET",
            headers: {
                Token: "kkajka",
            }
        })
        .then(resp => resp.json())
        .then(resp => {
            let task = resp[0];
            
            let head = document.createElement("h2");
            head.innerText = `задача #${id} (тип: ${task["type"]})`;
            res.append(head);

            return task["text"];
        })

        console.log(file);

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

        return res;
    }
}