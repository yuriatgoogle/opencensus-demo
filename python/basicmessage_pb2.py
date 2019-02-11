# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: basicmessage.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='basicmessage.proto',
  package='basicMessage',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x12\x62\x61sicmessage.proto\x12\x0c\x62\x61sicMessage\"\x1e\n\x0e\x62\x61\x63kendRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"\x1f\n\x0c\x62\x61\x63kendReply\x12\x0f\n\x07message\x18\x01 \x01(\t2T\n\x07\x62\x61\x63kend\x12I\n\x0bsendMessage\x12\x1c.basicMessage.backendRequest\x1a\x1a.basicMessage.backendReply\"\x00\x62\x06proto3')
)




_BACKENDREQUEST = _descriptor.Descriptor(
  name='backendRequest',
  full_name='basicMessage.backendRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='basicMessage.backendRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=36,
  serialized_end=66,
)


_BACKENDREPLY = _descriptor.Descriptor(
  name='backendReply',
  full_name='basicMessage.backendReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='message', full_name='basicMessage.backendReply.message', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=68,
  serialized_end=99,
)

DESCRIPTOR.message_types_by_name['backendRequest'] = _BACKENDREQUEST
DESCRIPTOR.message_types_by_name['backendReply'] = _BACKENDREPLY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

backendRequest = _reflection.GeneratedProtocolMessageType('backendRequest', (_message.Message,), dict(
  DESCRIPTOR = _BACKENDREQUEST,
  __module__ = 'basicmessage_pb2'
  # @@protoc_insertion_point(class_scope:basicMessage.backendRequest)
  ))
_sym_db.RegisterMessage(backendRequest)

backendReply = _reflection.GeneratedProtocolMessageType('backendReply', (_message.Message,), dict(
  DESCRIPTOR = _BACKENDREPLY,
  __module__ = 'basicmessage_pb2'
  # @@protoc_insertion_point(class_scope:basicMessage.backendReply)
  ))
_sym_db.RegisterMessage(backendReply)



_BACKEND = _descriptor.ServiceDescriptor(
  name='backend',
  full_name='basicMessage.backend',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=101,
  serialized_end=185,
  methods=[
  _descriptor.MethodDescriptor(
    name='sendMessage',
    full_name='basicMessage.backend.sendMessage',
    index=0,
    containing_service=None,
    input_type=_BACKENDREQUEST,
    output_type=_BACKENDREPLY,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_BACKEND)

DESCRIPTOR.services_by_name['backend'] = _BACKEND

# @@protoc_insertion_point(module_scope)
