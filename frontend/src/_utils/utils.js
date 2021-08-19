/**
 * @param {Array} arr 
 */
const arrayToCSV = ( arr ) =>{

    if (arr.length === 0)
        return '';

    let buff = '';

    for( let i=0; i<arr.length; ++i ){
        buff = arr[i].toString();
        if ( i < arr.length-1)
            buff += ',';
    }

    return buff;
};

export {
    arrayToCSV
}