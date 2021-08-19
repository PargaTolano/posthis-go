const fileToBase64 = ( file )=>new Promise((res, rej)=>{
    var reader  = new FileReader();

    reader.onloadend    = (     )  => res( reader.result );
    reader.onerror      = ( err )  => rej( err );
    reader.onabort      = ( err )  => rej( err );

    if ( file ) {
        reader.readAsDataURL(file);
    }else{
        rej();
    }
});


export {
    fileToBase64
}