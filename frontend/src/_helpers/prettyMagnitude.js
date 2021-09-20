import { roundToOneDecimal } from "./round";

function inRange(number, min, max) {
    return ( number >= min && number <= max);
}

export function prettyMagnitude( count ){

    if ( count === 0 ) return 0;
			
	if ( isNaN(count) || count === null || count === undefined)
		return null;
        
	return (
        ( count / 1e2  < 10 ) && count                  ||
        ( inRange (count /  1e5, 1, 10) ) && ( roundToOneDecimal( count / 1e3 ) + "K" )    ||
        ( inRange (count /  1e8, 1, 10) ) && ( roundToOneDecimal( count / 1e6 ) + "M"  )   ||
        ( inRange (count / 1e11, 1, 10) ) && ( roundToOneDecimal( count / 1e9 ) + "B"  )
    );
}