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
        "raid_resync_started_succ": 'RAID resync started',
        "raid_device_updated_succ": 'RAID device status reloaded',
        "raid_reassemble_started_succ": 'RAID config reloaded',
        "raid_device_deleted_succ": 'RAID device deleted',
        "raid_device_deleted_fail": 'RAID device delete failed',
        "raid_device_created_succ": 'RAID device created',
        "raid_device_created_fail": 'RAID device create failed',
    },
    'zh': {
        'disk_info_refreshed': '磁碟資訊已重新載入',
        "raid_resync_started_succ": 'RAID 重建已成功啟動',
        "raid_device_updated_succ": 'RAID 裝置狀態已重新載入',
        "raid_reassemble_started_succ": 'RAID 配置已重新載入',
        "raid_device_deleted_succ": 'RAID 裝置已刪除',
        "raid_device_deleted_fail": 'RAID 裝置刪除失敗',
        "raid_device_created_succ": 'RAID 裝置已建立',
        "raid_device_created_fail": 'RAID 裝置建立失敗',
    }
};