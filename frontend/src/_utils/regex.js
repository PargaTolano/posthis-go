const passwordRegex =
     /^[a-z][a-z0-9_.-]*/i;

const userNameRegex = /^[0-9A-Za-z_.-]+$/;

const emailRegex = 
     /([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/;

const tagRegex = /^[a-z0-9]+$/i;

export{
     passwordRegex,
     userNameRegex,
     emailRegex,
     tagRegex
}