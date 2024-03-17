BEGIN;

ALTER table films ALTER column released SET NOT NULL;
ALTER table films ALTER column description SET NOT NULL;
ALTER table films ALTER column rating SET NOT NULL;

ALTER table actors ALTER column birthday SET NOT NULL;
ALTER table actors ALTER column gender SET NOT NULL;

END;