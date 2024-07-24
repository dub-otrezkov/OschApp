
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
                // res += `
                //     <a class="tasklink" id="${resp[i]["id"]}" href="/tasks/${resp[i]["id"]}">задача #${resp[i]["id"]}</a>
                // `;
            }

            console.log(res);

            return res;

        })

        .catch(err => console.log(err));

        return res;
    }
}