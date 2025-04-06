/*
    Localization
*/
let languages = ['en', 'zh'];
let languageNames = {
    'en': 'English',
    'zh': '中文（正體）'
};
let currentLanguage = 'en';

//Initialize the i18n dom library
var i18n = domI18n({
    selector: '[i18n]',
    separator: ' // ',
    languages: languages,
    defaultLanguage: 'en'
});

$(document).ready(function(){
    let userLang = navigator.language || navigator.userLanguage;
    console.log("User language: " + userLang);
    userLang = userLang.split("-")[0];
    if (!languages.includes(userLang)) {
        userLang = 'en';
    }
    i18n.changeLanguage(userLang);
    currentLanguage = userLang;
});

// Update language on newly loaded content
function relocale(){
    i18n.changeLanguage(currentLanguage);
}

function setCurrentLanguage(newLanguage){
    let languageName = languageNames[newLanguage];
    currentLanguage = newLanguage;
    $("#currentLanguage").html(languageName);
    i18n.changeLanguage(newLanguage);
}

/* Other Translated messages */
function i18nc(key, language=undefined){
   if (language === undefined){
       language = currentLanguage;
    }

    let translatedMessage = translatedMessages[language][key];
    if (translatedMessage === undefined){
        translatedMessage = translatedMessages['en'][key];
    }
    if (translatedMessage === undefined){
        translatedMessage = key;
    }
    return translatedMessage;
}

let translatedMessages = {
    'en': {
        'disk_info_refreshed': 'Disk information reloaded',
    },
    'zh': {
        'disk_info_refreshed': '磁碟資訊已重新載入',
    }
};