
function interpreter(query) {

    let sentenceSplit = query.split(".")
    
    for(let i = 0; i < sentenceSplit.length; i++) {

        sentenceSplit[i] = sentenceSplit[i].trim()

    }
    
    sentenceSplit = sentenceSplit.filter(s => s != '')

    console.log("Sentence Split Array:")
    console.log(sentenceSplit)

    processQueries(sentenceSplit)

}

function processQueries(queryArr) {

    for(let query of queryArr) {
        
        console.log(query)

    }

}


function isMathQuery(query) {

    let mathQuery = query.split(" ")
    mathQuery = mathQuery.filter(s => s != '')
    console.log(mathQuery)

    let numCount = 0 
    for(let mathQ of mathQuery) {

        if(!isNaN(parseInt(mathQ))) numCount++

    }

    return numCount > 0 ? true : false 

}

console.log(isMathQuery("Hello he haa"))

test = "Hello world. This is a program and it is pretty good. I like it quite a bit a bit. a."

// interpreter(test)