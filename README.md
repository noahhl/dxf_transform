Move DXFs to have their lower left corner at 0,0; scale them; and close polylines.

Usage:
```
dxf_transform -in=MI.dxf -out=MI-adjust.dxf -translate=true -close=true -scale=1.5   
```

Works on DXFs that are entirely polylines. For everything else, I dunno.
