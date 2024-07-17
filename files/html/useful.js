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
