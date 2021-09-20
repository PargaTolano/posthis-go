import { toast } from 'react-toastify';

function makeToast( content, type ){
    toast[type]( content, {
        position: 'top-left',
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined,
    });
}

export function makeInfoToast( content ){
    makeToast(content, 'info');
}

export function makeSuccessToast( content ){
    makeToast(content, 'success');

}

export function makeWarningToast( content ){
    makeToast(content, 'warn');

}

export function makeErrorToast( content, ){
    makeToast(content, 'error');
}