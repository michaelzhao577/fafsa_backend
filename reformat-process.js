const puppeteer = require('puppeteer');
const login = require('./login.json');
const http = require('http');
const { resolve } = require('path');

// Promisified HTTP request to GO backend to request the answer for a specific question/
// If HTTP request returns an answer, the promise resolves to said answer.
// Will throw an error if no answer is returned.
// Question: Should we query one answer at a time or just all of the answers for an user?
function requestData(id, question){
    let address = "http://localhost:8080/data/"+id+"/"+question;
    return new Promise ((resolve, reject) => {
        http.get(address, (response) => {
            let answer = "";
            response.on("data", (data) => {
                answer = JSON.parse(data);
            });
            response.on("end", () => {
                if (answer == ""){
                    reject("");
                }
                else{
                    resolve(answer);
                }
            });
        });
    });
}


// New reusable async function to look for selector
function lookForSelector(currPage, selector){
    return new Promise((resolve, reject) => {
        try {
            // Timeout is 10 ms since we are only looking for selector once page is laoded
            currPage.waitForSelector(selector, {timeout:10}).then((value)=>{
                if (value === null){
                    // If element is visibility:hidden it will resolve to null
                    console.log("NULL");
                    resolve(false);
                }else{
                    // We are able to retrieve DOM element
                    console.log("DOM element found:"+ value);
                    resolve(true);
                }
            }, (error) => {
                // Unable to find DOM element/time out
                console.log("Promise rejected:" +error);
            })
        } catch (error){
            // Potential throws an error, so catch
            console.log("Error: ", error);
        }
    });
}

// Similar to method above, but we attach typing in a value 
async function typeInSelector(page, selector, answer){
    try {
        await page.type(selector, answer);
    } catch (error){
        console.log("Error: ", error);

    }
}

// Michael
async function typeInSelector(currPage, selector, text){
    try {
      currPage.waitForSelector(selector)
        .then((value) => {
          currPage.type(value, text)
        })
        .catch((error) => {
        console.log("Could not type in selector. error:", error)
        sendErrorMessage("could not type in selector " + selector + ". " + error)
        })
    } catch (error) {
        console.log("Error:", error);
    }
  }

// Michael
async function clickSamePage(currPage, selector) {
    try {
      currPage.waitForSelector(selector, {  timeout: 5000 })
        .then(() => {
          currPage.click(selector)
          console.log("click successful"); 
        })
        .catch((error) => {
          console.log("Promise rejected. Could not find selector. Error:", error)
          sendErrorMessage("could not find selector " + selector + ". " + error)
        }
      )
    } catch (error) {}
      console.log("error:", error)
  }

// Michael
async function clickToNewPage(currPage, selector) {
    try {
      currPage.waitForNavigation()
        .then(() => {
          currPage.click(selector)
        })
        .catch((error) => {
          console.log("Promise rejected. Could not click to new page. Error:", error)
          sendErrorMessage("Failed to click to next page. could not find selector " + selector + ". " + error)
        })
    } catch (error) {
      console.log("Promise rejected with error:", error)
    }
  }

async function sendErrorMessage(message) {
    // create reusable transporter object using the default SMTP transport
    let transporter = nodemailer.createTransport({
      host: 'outlook.office365.com',
      port: 587,
      secure: false, // true for 465, false for other ports
      auth: {
        user: "mifdhsaiofhkasjdbjgasdf@outlook.com",
        pass: "EDUrain!"
      },
    });
  
    // send mail with defined transport object
    let info = await transporter.sendMail({
      from: '"mjz" <mifdhsaiofhkasjdbjgasdf@outlook.com>', // sender address
      to: "michael_zhao@yahoo.com", // list of receivers
      subject: "fafsa submitter error", // Subject line
      text: message, // plain text body
    });
  
    console.log("Message sent: %s", info.messageId);
    // Message sent: <b658f8ca-6296-ccf4-8306-87d57a0b4321@example.com>
  
  }
  

(async () => {
  const browser = await puppeteer.launch({
    devtools: true,
    headless: false
  });
  // ===========================================================================

  const page = await browser.newPage();
  // Iterating through each HTML element in a single JSON and taking action
  await page.goto('https://fafsa.ed.gov/spa/fafsa/#/LOGIN;direction=next', { waitUntil: 'networkidle2' });
  await page.waitForSelector('input[id="login-tab-group:1"]');
  await page.click('input[id="login-tab-group:1"]', { delay: 100 });
  await page.waitForSelector('#pii-tab-content');
  for (var element in login){
      const found = await lookForSelector(page, element);
      if (found == true){
        const data = await requestData(1, login[element].question);
        await typeInSelector(page, element, data);
      } else{

      }
  }
})();

