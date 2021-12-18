const puppeteer = require('puppeteer');
const fs = require('fs');

const URL = 'https://intra.assistants.epita.fr/';
const DOWLOAD_PATH = '/output';

function delay(time) {
    return new Promise(function (resolve) {
        setTimeout(resolve, time)
    });
}

async function scanLinks(page, linkType) {
    return await page.evaluate(async (linkType) => {
        const r = [];
        const links = document.getElementsByTagName('a');
        for (let i = 0; i < links.length; i++) {
            if (links[i].href.includes(linkType)) {
                r.push(links[i].href);
            }
        }
        return r;
    }, linkType);
}

function storeFileURL(url, path) {
    console.log('Downloading file...', url, path);
    fs.writeFile(DOWLOAD_PATH + '/' + path + '.url', url, function (err) {
        if (err) {
            return console.log(err);
        }
        console.log("The file was saved!");
    });
}

async function archive() {
    const browser = await puppeteer.launch({
        headless: false, // does not work in headless mode
        userDataDir: './userData',
        executablePath: process.env.CHROME_BIN,
        args: ['--no-sandbox', '--disable-gpu', '--disable-dev-shm-usage']
    });
    const page = await browser.newPage();
    await page._client.send('Page.setDownloadBehavior', { behavior: 'allow', downloadPath: DOWLOAD_PATH });

    await page.goto(URL);
    await delay(2000);

    if (page.url().includes('login')) {
        console.log('Logging in...');
        await page.evaluate((p) => {
            const login = document.querySelector('#id_username');
            login.value = p;
        }, process.env.LOGIN);

        await page.evaluate((p) => {
            const pass = document.querySelector('#id_password');
            pass.value = p;
            pass.focus();
        }, process.env.PASSWORD);

        await delay(100);
        await page.keyboard.press('Enter');
        await delay(3000);
    }
    console.log('Logged in.');


    console.log('Fetching activities...');
    const links = await scanLinks(page, 'activity');
    console.log('Fetched activities.');

    console.log('Found ' + links.length + ' activities.');
    await page.exposeFunction('storeFileURL', storeFileURL);
    for (let i = 0; i < links.length; i++) {
        await page.goto(links[i]);
        await delay(3000);

        console.log('Fetching assignments...');
        const assignments = await scanLinks(page, 'assignment');
        console.log('Fetched assignments.');

        for (let j = 0; j < assignments.length; j++) {
            await page.goto(assignments[j]);
            await delay(3000);

            await page.evaluate(async () => {

                function delay(time) {
                    return new Promise(function (resolve) {
                        setTimeout(resolve, time)
                    });
                }

                const links = document.querySelectorAll('[theme="small icon primary"]');

                for (let i = 0; i < links.length; i++) {
                    if (links[i].firstChild.icon.includes('file')) {
                        console.log('Downloading file...');
                        links[i].click();
                        
                        await window.storeFileURL(window.location.href, links[i].parentNode.parentNode.nextSibling.textContent);
                        await delay(1000);
                    }
                }
            });
        }
    }

    await browser.close();
}

archive();