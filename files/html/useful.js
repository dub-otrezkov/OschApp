async function combinePromises(...funcs) {
    const results = await Promise.all(funcs.map(func => func()));
    console.log(results)

    var vls = {}
    results.forEach((result, index) => {
        vls[funcs[index].name] = result;
    });

    console.log(vls);
    return vls;
}


var tasksPages = {
    oks: null,

    displaytasks: async function(document, apiUrl, UserId) {
        async function getTasks() {
            return fetch(apiUrl).then(
                response => {
                    if (response.status != 200) return [];
                    else return response.json();
                }
            )
        };

        async function getSubmissions() {

            if (this.oks != null) return [];

            return fetch(`/api/get/Submissions?user=${UserId}&status=1`).then(
                response => {
                    if (response.status != 200) return [];
                    else return response.json();
                }
            )
        };

        combinePromises(getTasks, getSubmissions).then(res => {

            console.log(res);
            
            var linksContainer = document.getElementById("linksContainer");

            console.log(linksContainer);
            linksContainer.innerHTML = "";

            console.log(res);
            console.log(res["getTasks"]);

            let tasks = res["getTasks"];
            let subs = res["getSubmissions"];


            console.log("subs: ", subs);

            if (this.oks == null && subs != null) {
                this.oks = {};
                subs.forEach(sub => {
                    console.log(sub);
                    this.oks[sub["task_id"]] = 1;
                });
            }
            
            console.log("oks: ", this.oks);

            // var j = getSubmissions();
            // console.log("gs: ", j);

            tasks.forEach(name => {
                const link = document.createElement('a');
                
                link.textContent = `#${name["id"]} (тип ${name["type"]})`;
                link.className = `tasklink`;
                link.href = `/tasks/${name["id"]}`;

                if (this.oks[name["id"]] == 1) {
                    link.id = "ok_task";
                }

                linksContainer.appendChild(link);
            });
        });
    },

    Task: {
        cor_ans: "",

        getTask: async function (document, apiUrl) {
            fetch(apiUrl)
            .then(response => {
                if (!response.ok) return [];
                return response.json();
            })
            .then(async data => {

                var tsk = data[0];
                console.log(tsk);
                var cont = document.getElementById("content");

                var header = document.getElementById("name");
                header.innerHTML = `задача #${tsk["id"]} (тип ${tsk["type"]})`;

                this.cor_ans = tsk["ans"];

                console.log("parsed: ", this.cor_ans);

                const res = await fetch(`/api/files/content/${tsk.text}`);
                const html = await res.text();
                console.log(html);
                cont.innerHTML = `${html} <br> <div id="subs"></div>`;
            })
            .then(() => {
                console.log("ans: ", this.cor_ans);
            })

            
        }
    }
}