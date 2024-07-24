
var taskApp = {
    getTasksList: async function(params, elemid) {
        var res = "";

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
                res += `
                    <a class="tasklink" id="${resp[i]["id"]}" href="/tasks/${resp[i]["id"]}">задача #${resp[i]["id"]}</a>
                `;
            }

            console.log(res);

            return res;

        })

        .catch(err => console.log(err));

        return res;
    }
}