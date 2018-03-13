# GoBufrKit

> An awesome project.

* 201, 202, 203 must be cancelled or completed (?) within a replication session
    - This check is not implemented

* 201, 202, 203 Must not be nested with 207
    - This check is not implemented

* 203YYY new refval
    - Technically it is possible for each Compressed subset has different values for
      new refval definition. However, is this really practical? It does not seem
      to be useful and create much trouble for programming.
    - The descriptors follow 203YYY can in theory be more than just Element descriptor.
      However it does not make sense to have non-element descriptor. This should be
      a lint error

* 204YYY associated fields
    - Does it apply to skipped local descriptors, operator for inserting strings etc.?
        * NOT going to apply?
    - In theory, associated fields can be applied to the target element descriptors
      of an assessment session. This creates much problem for implementation for backref
      counting.


* 221YYY data not present
    - How to count the YYY descriptors? Count only top level descriptors,
      similar to how replication is counted?
    - Do associated fields apply to those NOT present descriptors?
        * NOT going to apply
    - Should it apply to skip local descriptor? i.e. should local descriptor still
      be skipped by reading YYY bits or no action at all (as not present)?
        * Decide to let both rules take effect

* 235000 cancel backward data reference
    - This implies that multiple bitmaps may be applied to the same section of target descriptors?
    - If a new bitmap wants to be applied to a different descriptor section, 235000 has to be used

* 236000 define reusable bitmap
    - Is it not allowed to have standalone 236000 session, i.e. a session that is not
      inside a parent assessment session.
        * Technically it is possible
    - Bitmap are applied to to flattened element descriptors, NOT just the top level
      ones, e.g. (amv2_87.bufr)
        * 303250
            - 002252
            - 002023
            - 007004
            - 011001
            - 011002
            - 002197
            - 002198
            - 012193
        * 222000
        * 236000  (this bitmap definition applies to the element descriptors of 303250)
        * 101103
        * 031031

* In theory, the sandwiched descriptors could be sequence or replication descriptors,
  NOT just element descriptors.


From Manual on Codes
* Where an operator descriptor requires a data present bit-map of length N to complete
  the operator definition, the N consecutive element descriptors which correspond to
  the N data entities to which the N bit values refer shall end with the element
  descriptor which immediately precedes the first such operator, or with the element
  descriptor which immediately precedes the first occurrence of such an operator
  following the occurrence of a cancel backward reference operator.
    - This basically says that the sequence of target element descriptors are defined once for
      all subsequent operator descriptors. If a new sequence of target element descriptors
      need to be defined, a cancel backward reference operator must be used first.
      0XXYYY  sequence of target element descriptor
      0XXYYY
      0XXYYY
      224000  target is the first 3 element descriptor
      ......
      ......
      223000  target is still the first 3 element descriptor
      ......
      ......
      235000  cancel backward reference
      0XXYYY  new sequence of target element descriptor
      0XXYYY
      0XXYYY
      0XXYYY
      222000  target is now the new sequence of element descriptor
      ......
      ......


- 207003.bufr
    * Compressed
    * 2 subsets
    * 201
    * 202
    * 207 YYY

- uegabe.bufr
    * 1 subset
    * Associated field

- b002_95.bufr
    * 1 subset
    * Skipped local descriptor

- b005_89.bufr
    * Compressed
    * 128 subsets
    * Quality info, bitmap define, recall first-order stats

- g2nd_208.bufr
    * Compressed
    * 18 subsets
    * 207
    * 224, bitmap define

- amv2_87.bufr
    * Compressed
    * 128 subsets
    * Quality info, bitmap define, recall

- asr3_190.bufr
    * Compressed
    * 128 subsets
    * Quality info, bitmap define, recall, first-order stats

- jaso_214.bufr
    * Compressed
    * 128 subsets
    * Associated fields