import GetCookie from '../../scripts/general.js';

var UserStats = {
    getStats: async function () {
        let res = document.createElement("div");

        let userId = GetCookie("userId")

        await fetch(`/api/stats/${userId}`, {
            method: "GET",
        })
        .then(resp => {
            if (resp.ok) return resp.json();
            else return [];
        })
        .then(resp => {
            console.log(resp);

            resp.forEach(exam => {
                let cont = document.createElement("div");

                let head = document.createElement("a");
                head.text = `Экзамен #${exam["exam_id"]}`;
                head.href = `/exams/${exam["exam_id"]}`;
                let d = document.createElement("table");
                let up = document.createElement("tr");
                let dn = document.createElement("tr");

                
                for (let k in exam["ans"]) {
                    
                    {
                        let p = document.createElement("td");
                        let a = document.createElement("a");
                        a.innerText = k;
                        a.href = `/tasks/${k}`;
                        p.append(a);

                        up.append(p);
                    }

                    {
                        let p = document.createElement("td");

                        switch(exam["ans"][k]) {
                            case 0:
                                p.innerText = "0";
                                p.classList += " wa_task";
                                break;
                            case 1:
                                p.innerText = "1";
                                p.classList += " ok_task";
                                break;
                            case 2:
                                p.innerText = "-";
                                p.classList += " nt_task";
                        }

                        dn.append(p);
                    }
                }

                d.append(up, dn);

                cont.append(head, d);
                res.append(cont)
            });
        })

        return res;
    }
}

export default UserStats